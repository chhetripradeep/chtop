package chtop

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chhetripradeep/chtop/pkg/metric"
	"github.com/spf13/viper"

	"github.com/chhetripradeep/chtop/pkg/model"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

func Run(url string) error {
	currentTheme, err := theme.LoadViper(viper.GetViper(), true)
	if err != nil {
		return err
	}

	m := model.Model{
		Endpoint: url,
		Metrics: metric.Metrics{
			metric.NewMetric("ClickHouseProfileEvents_Query"),
			metric.NewMetric("ClickHouseProfileEvents_SelectQuery"),
			metric.NewMetric("ClickHouseProfileEvents_InsertQuery"),
			metric.NewMetric("ClickHouseMetrics_PartsActive"),
			metric.NewMetric("ClickHouseMetrics_TCPConnection"),
		},
		Theme: currentTheme,
	}

	program := tea.NewProgram(m, tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		return err
	}
	return nil
}
