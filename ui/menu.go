package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	cursor int
}

func NewMenuModel() MenuModel {
	return MenuModel{
		cursor: 0,
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
			selectedOption := menuOptions[m.cursor]

			return m, func() tea.Msg {
				return navigateToMsg{view: selectedOption.View}
			}
		case key.Matches(msg, keyMap.MoveUp):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(menuOptions) - 1
			}
		case key.Matches(msg, keyMap.MoveDown):
			m.cursor++
			if m.cursor >= len(menuOptions) {
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Welcome to gowizcli! A Wiz lights client written in Go. Select an option:\n\n")

	for i, o := range menuOptions {
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

var menuOptions = []struct {
	Name string
	View ViewType
}{
	{Name: "Discover lights in local network", View: ViewDiscover},
	{Name: "Show current lights", View: ViewShow},
	{Name: "Erase current lights", View: ViewEraseAll},
}
