package chtop

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhetripradeep/chtop/pkg/query"
	"github.com/spf13/viper"

	"github.com/chhetripradeep/chtop/pkg/metric"
	"github.com/chhetripradeep/chtop/pkg/model"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

func Run(metricsUrl, queriesUrl, database, username, password string) error {
	themes, err := theme.LoadTheme(viper.GetViper(), true)
	if err != nil {
		return err
	}

	metrics, err := metric.LoadMetrics(viper.GetViper(), true)
	if err != nil {
		return err
	}

	queries, err := query.LoadQueries(viper.GetViper(), true)
	if err != nil {
		return err
	}

	m := model.Model{
		Theme:              themes,
		MetricsEndpoint:    metricsUrl,
		QueriesEndpoint:    queriesUrl,
		ClickHouseDatabase: database,
		ClickHouseUsername: username,
		ClickHousePassword: password,
		ClickHouseMetrics:  metrics,
		ClickHouseQueries:  queries,
	}

	program := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so that we can track the mouse wheel
	)
	_, err = program.Run()
	if err != nil {
		return err
	}
	return nil
}
