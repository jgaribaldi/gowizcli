package ui

import (
	"gowizcli/client"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
		switch {
		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, keyMap.Select):
			m.selected = m.cursor
			selectedOption := m.options[m.selected]

			return m, func() tea.Msg {
				view := viewCommandMap[selectedOption.Type]
				return navigateToMsg{view: view}
			}
		case key.Matches(msg, keyMap.MoveUp):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		case key.Matches(msg, keyMap.MoveDown):
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
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

var keyMap = struct {
	Quit     key.Binding
	Select   key.Binding
	MoveUp   key.Binding
	MoveDown key.Binding
}{
	Quit:     key.NewBinding(key.WithKeys("ctrl+c", "q", "esc"), key.WithHelp("ctrl+c / q / esc", "quit program")),
	Select:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	MoveUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up arrow / k", "move up")),
	MoveDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down arrow / j", "move down")),
}
