package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"

	"github.com/chhetripradeep/chtop/pkg/client"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

const (
	defaultFps = time.Duration(30)
)

type Metric struct {
	Name       string
	Latest     float64
	Datapoints []float64
}

func NewMetric(name string) Metric {
	return Metric{
		Name: name,
	}
}

func (m *Metric) Update(value float64) {
	m.Latest = value
	m.Datapoints = append(m.Datapoints, value)
}

var metrics = []Metric{
	NewMetric("ClickHouseProfileEvents_Query"),
	NewMetric("ClickHouseProfileEvents_SelectQuery"),
	NewMetric("ClickHouseProfileEvents_InsertQuery"),
	NewMetric("ClickHouseMetrics_PartsActive"),
	NewMetric("ClickHouseMetrics_TCPConnection"),
}

type Model struct {
	Client *client.Client
	Theme  *theme.Theme
}

// Init inits the bubbletea model for use
func (m Model) Init() tea.Cmd {
	return nil
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
	default:
		return m, tick()
	}
}

// View shows the current state of the chtop
func (m Model) View() string {
	var plot string
	for i := range metrics {
		value, _ := m.Client.GetMetric(metrics[i].Name)
		floatValue, _ := strconv.ParseFloat(strings.TrimSpace(*value), 64)
		metrics[i].Update(floatValue)
		graph := asciigraph.Plot(
			metrics[i].Datapoints,
			asciigraph.Height(m.Theme.Graph.Height),
			asciigraph.Width(m.Theme.Graph.Width),
			asciigraph.Precision(m.Theme.Graph.Precision),
			asciigraph.SeriesColors(m.Theme.GraphColor()),
		)
		plot += fmt.Sprintf("\n\n")
		plot += fmt.Sprintf(setTitle(metrics[i].Name).String())
		plot += fmt.Sprintf("\n\n%s\n\n", graph)
		plot += fmt.Sprintf(setFooter(fmt.Sprintf("Current Value: %.2f", metrics[i].Latest)).String())
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
