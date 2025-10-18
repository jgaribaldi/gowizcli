package ui

import (
	"fmt"
	"gowizcli/client"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	cursor   int
	options  []client.Option
	selected int
}

func NewMenuModel() MenuModel {
	var options []client.Option
	options = make([]client.Option, 0)

	for _, o := range client.Options {
		options = append(options, o)
	}

	return MenuModel{
		cursor:  0,
		options: options,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.selected = m.cursor
			fmt.Println("blah")
			selectedOption := m.options[m.selected]
			fmt.Printf("Selected option is %v\n", selectedOption)
			return m, func() tea.Msg {
				view := viewCommandMap[selectedOption.Type]
				return navigateToMsg{view: view}
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Welcome to gowizcli! A Wiz lights client written in Go. Select an option:\n\n")

	for i, o := range m.options {
		if m.cursor == i {
			s.WriteString("[*] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(o.Name)
		s.WriteString("\n")
	}

	return s.String()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
