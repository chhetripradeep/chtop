package theme

import (
	"github.com/guptarohit/asciigraph"
)

type Theme struct {
	Border string `json:"border" toml:"border" yaml:"border"`
	File   string `json:"file" toml:"file" yaml:"file"`
	Graph  Graph  `json:"graph" toml:"graph" yaml:"graph"`
}

type Graph struct {
	Color     string `json:"color" toml:"color" yaml:"color"`
	Height    int    `json:"height" toml:"height" yaml:"height"`
	Width     int    `json:"width" toml:"width" yaml:"width"`
	Precision uint   `json:"precision" toml:"precision" yaml:"precision"`
}

func DefaultTheme() *Theme {
	return &Theme{
		Border: "blue",
		Graph: Graph{
			Color:     "blue",
			Height:    5,
			Width:     60,
			Precision: 2,
		},
	}
}

func (t Theme) GraphColor() asciigraph.AnsiColor {
	graphColor := asciigraph.Blue
	if color, ok := asciigraph.ColorNames[t.Graph.Color]; ok {
		graphColor = color
	}
	return graphColor
}
