package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().
	Margin(0, 2)

var titleStyle = lipgloss.NewStyle()

var helplineStyle = lipgloss.NewStyle()

var tableStyle2 = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

var boxStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder())

func tableStyle() lipgloss.Style {
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Align(lipgloss.Center)
}
