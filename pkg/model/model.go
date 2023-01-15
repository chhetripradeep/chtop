package model

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"

	"github.com/chhetripradeep/chtop/pkg/metric"
	"github.com/chhetripradeep/chtop/pkg/query"
	"github.com/chhetripradeep/chtop/pkg/theme"
	"github.com/chhetripradeep/chtop/pkg/utils"
)

const (
	defaultFps                 = time.Duration(30)
	useHighPerformanceRenderer = false
)

type Model struct {
	Error              error
	Ready              bool
	Theme              *theme.Theme
	Spinner            spinner.Model
	Viewport           viewport.Model
	MetricsEndpoint    string
	QueriesEndpoint    string
	ClickHouseDatabase string
	ClickHouseUsername string
	ClickHousePassword string
	ClickHouseMetrics  metric.ClickHouseMetrics
	ClickHouseQueries  query.ClickHouseQueries
}

type statusMsg int

type dataMsg int

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// checkUrl checks the http endpoint
func checkUrl(url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(url)
		if err != nil {
			return errMsg{err}
		}
		defer resp.Body.Close()
		return statusMsg(resp.StatusCode)
	}
}

// checkEndpoint checks the tcp endpoint
func checkEndpoint(endpoint string) tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("tcp", endpoint)
		if err != nil {
			return errMsg{err}
		}
		defer conn.Close()
		return statusMsg(200)
	}
}

// Init inits the bubbletea model for use
func (m Model) Init() tea.Cmd {
	if m.MetricsEndpoint != "" {
		return checkUrl(m.MetricsEndpoint)
	}
	if m.QueriesEndpoint != "" {
		return checkEndpoint(m.QueriesEndpoint)
	}
	return nil
}

// Query hits the prometheus exporter endpoint of clickhouse
// and returns string containing the metric name and value
func (m Model) Query(metric string) (*string, error) {
	resp, err := http.Get(m.MetricsEndpoint)
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

// Execute runs a clickhouse query and returns the response
func (m Model) Execute(sql string) (*string, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{m.QueriesEndpoint},
		Auth: clickhouse.Auth{
			Database: m.ClickHouseDatabase,
			Username: m.ClickHouseUsername,
			Password: m.ClickHousePassword,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 30,
		},
		DialTimeout: 5 * time.Second,
	})
	defer conn.Close()

	rows, err := conn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var val string
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			return nil, err
		}
	}
	return &val, nil
}

// Update updates the bubbletea model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
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
	case tea.WindowSizeMsg:
		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height)
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height
		}
		hex := utils.ColorNameToHex(m.Theme.Border)
		m.Viewport.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(hex)).
			Padding(2, 2, 2, 2)
		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
	case errMsg:
		m.Error = msg
		return m, tea.Quit
	case dataMsg:
		m.Ready = true
		return m, m.UpdateData()
	default:
		return m, m.UpdateData()
	}
	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// UpdateMetricsData updates the data from prometheus exporter's metrics endpoint
func (m Model) UpdateMetricsData() {
	for i := range m.ClickHouseMetrics.Metrics {
		value, err := m.Query(m.ClickHouseMetrics.Metrics[i].Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to query endpoint:", m.MetricsEndpoint, "for metric:", m.ClickHouseMetrics.Metrics[i].Name)
			continue
		}

		floatValue, err := strconv.ParseFloat(strings.TrimSpace(*value), 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to parse value for metric:", m.ClickHouseMetrics.Metrics[i].Name)
			continue
		}
		m.ClickHouseMetrics.Metrics[i].Update(floatValue)
	}
}

// UpdateQueriesData updates the data from clickhouse queries output
func (m Model) UpdateQueriesData() {
	for i := range m.ClickHouseQueries.Queries {
		value, err := m.Execute(m.ClickHouseQueries.Queries[i].Sql)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to execute query:", m.ClickHouseQueries.Queries[i].Name, "at endpoint:", m.QueriesEndpoint)
		}

		floatValue, err := strconv.ParseFloat(strings.TrimSpace(*value), 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to parse value for metric:", m.ClickHouseQueries.Queries[i].Name)
			continue
		}
		m.ClickHouseQueries.Queries[i].Update(floatValue)
	}
}

// UpdateData updates the complete data
func (m Model) UpdateData() tea.Cmd {
	return tea.Tick(time.Minute/defaultFps, func(t time.Time) tea.Msg {
		if m.MetricsEndpoint != "" {
			m.UpdateMetricsData()
		}
		if m.QueriesEndpoint != "" {
			m.UpdateQueriesData()
		}
		return dataMsg(1)
	})
}

func (m Model) ViewMetricsData() string {
	var plot string
	for i := range m.ClickHouseMetrics.Metrics {
		caption := m.ClickHouseMetrics.Metrics[i].Alias + fmt.Sprintf(" (Current Value: %.2f)\n\n", m.ClickHouseMetrics.Metrics[i].Latest)

		graph := asciigraph.Plot(
			m.ClickHouseMetrics.Metrics[i].Datapoints,
			asciigraph.Height(m.Theme.Graph.Height),
			asciigraph.Width(m.Theme.Graph.Width),
			asciigraph.Precision(m.Theme.Graph.Precision),
			asciigraph.SeriesColors(m.Theme.GraphColor()),
			asciigraph.Caption(caption),
		)

		plot += lipgloss.JoinVertical(
			lipgloss.Top,
			graph,
		)
		plot += "\n\n"
	}
	return plot
}

func (m Model) ViewQueriesData() string {
	var plot string
	for i := range m.ClickHouseQueries.Queries {
		caption := m.ClickHouseQueries.Queries[i].Name + fmt.Sprintf(" (Current Value: %.2f)\n\n", m.ClickHouseQueries.Queries[i].Latest)

		graph := asciigraph.Plot(
			m.ClickHouseQueries.Queries[i].Datapoints,
			asciigraph.Height(m.Theme.Graph.Height),
			asciigraph.Width(m.Theme.Graph.Width),
			asciigraph.Precision(m.Theme.Graph.Precision),
			asciigraph.SeriesColors(m.Theme.GraphColor()),
			asciigraph.Caption(caption),
		)

		plot += lipgloss.JoinVertical(
			lipgloss.Top,
			graph,
		)
		plot += "\n\n"
	}
	return plot
}

// View shows the current state of the complete data
func (m Model) View() string {
	var metricsPlot, queriesPlot, finalPlot string

	// If we don't have data to show in view
	if !m.Ready {
		return m.Spinner.View() + "  Initializing..."
	}

	if m.MetricsEndpoint != "" {
		metricsPlot = m.ViewMetricsData()
	}
	if m.QueriesEndpoint != "" {
		queriesPlot = m.ViewQueriesData()
	}

	finalPlot = lipgloss.JoinHorizontal(
		lipgloss.Left,
		metricsPlot,
		queriesPlot,
	)
	m.Viewport.SetContent(finalPlot)
	return m.Viewport.View()
}
