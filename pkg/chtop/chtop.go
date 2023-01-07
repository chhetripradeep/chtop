package chtop

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhetripradeep/chtop/pkg/query"
	"github.com/spf13/viper"

	"github.com/chhetripradeep/chtop/pkg/metric"
	"github.com/chhetripradeep/chtop/pkg/model"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

func Run(metricsUrl, queriesUrl string) error {
	currentTheme, err := theme.LoadTheme(viper.GetViper(), true)
	if err != nil {
		return err
	}

	currentMetrics, err := metric.LoadMetrics(viper.GetViper(), true)
	if err != nil {
		return err
	}

	currentQueries, err := query.LoadQueries(viper.GetViper(), true)
	if err != nil {
		return err
	}

	m := model.Model{
		MetricsEndpoint:   metricsUrl,
		QueriesEndpoint:   queriesUrl,
		ClickHouseMetrics: currentMetrics,
		ClickHouseQueries: currentQueries,
		Theme:             currentTheme,
	}

	program := tea.NewProgram(m, tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		return err
	}
	return nil
}
