package chtop

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"

	"github.com/chhetripradeep/chtop/pkg/client"
	"github.com/chhetripradeep/chtop/pkg/model"
	"github.com/chhetripradeep/chtop/pkg/theme"
)

func Run(url string) error {
	currentTheme, err := theme.LoadViper(viper.GetViper(), true)
	if err != nil {
		return err
	}

	m := model.Model{
		Theme: currentTheme,
		Client: &client.Client{
			Url: url,
		},
	}
	program := tea.NewProgram(m, tea.WithAltScreen())

	_, err = program.Run()
	if err != nil {
		return err
	}
	return nil
}
