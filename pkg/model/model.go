package model

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"

	"github.com/chhetripradeep/chtop/pkg/metric"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

const (
	defaultFps = time.Duration(30)
)

type Model struct {
	ClickHouseMetrics metric.ClickHouseMetrics
	Endpoint          string
	Theme             *theme.Theme
	Err               error
}

type statusMsg int

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

func check(url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(url)
		if err != nil {
			return errMsg{err}
		}
		defer resp.Body.Close()
		return statusMsg(resp.StatusCode)
	}
}

// Init inits the bubbletea model for use
func (m Model) Init() tea.Cmd {
	return check(m.Endpoint)
}

func (m Model) Query(metric string) (*string, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), metric) {
			val := strings.TrimLeft(scanner.Text(), metric)
			return &val, nil
		}
	}
	return nil, errors.New("unable to find the requested metric")
}

// Update updates the bubbletea model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			fallthrough
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.Err = msg
		return m, tea.Quit
	default:
		return m, tick()
	}
}

// View shows the current state of the chtop
func (m Model) View() string {
	var plot string
	for i := range m.ClickHouseMetrics.Metrics {
		value, err := m.Query(m.ClickHouseMetrics.Metrics[i].Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to query metrics endpoint:", m.Endpoint, "for metric:", m.ClickHouseMetrics.Metrics[i].Name)
			continue
		}

		floatValue, err := strconv.ParseFloat(strings.TrimSpace(*value), 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to parse value for metric:", m.ClickHouseMetrics.Metrics[i].Name)
			continue
		}

		m.ClickHouseMetrics.Metrics[i].Update(floatValue)
		graph := asciigraph.Plot(
			m.ClickHouseMetrics.Metrics[i].Datapoints,
			asciigraph.Height(m.Theme.Graph.Height),
			asciigraph.Width(m.Theme.Graph.Width),
			asciigraph.Precision(m.Theme.Graph.Precision),
			asciigraph.SeriesColors(m.Theme.GraphColor()),
		)

		plot += fmt.Sprintf("\n\n")
		plot += lipgloss.JoinVertical(
			lipgloss.Top,
			setTitle(m.ClickHouseMetrics.Metrics[i].Name).String(),
			graph,
			setFooter(fmt.Sprintf("Current Value: %.2f\n\n", m.ClickHouseMetrics.Metrics[i].Latest)).String(),
		)
	}
	return plot
}

func setTitle(text string) lipgloss.Style {
	return lipgloss.NewStyle().
		MarginLeft(20).
		MarginRight(5).
		Padding(0, 1).
		Italic(false).
		Bold(true).
		Border(lipgloss.NormalBorder(), true, true).
		SetString(text)
}

func setFooter(text string) lipgloss.Style {
	return lipgloss.NewStyle().
		MarginLeft(20).
		MarginRight(5).
		Padding(0, 1).
		Italic(false).
		SetString(text)
}

type tickMsg struct {
	Time time.Time
}

func tick() tea.Cmd {
	return tea.Tick(time.Minute/defaultFps, func(t time.Time) tea.Msg {
		return tickMsg{Time: t}
	})
}
