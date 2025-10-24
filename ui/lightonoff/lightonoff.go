package lightonoff

import (
	"gowizcli/client"
	"gowizcli/ui/common"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	client    *client.Client
	table     table.Model
	cmdStatus common.CommandStatus
}

func NewModel(c *client.Client) Model {
	columns := []table.Column{
		{Title: "ID", Width: 40},
		{Title: "MAC Address", Width: 20},
		{Title: "IP Address", Width: 20},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)

	return Model{
		client:    c,
		table:     t,
		cmdStatus: common.NewCommandStatus(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
		case tea.KeyShiftTab:
		case tea.KeyEnter:
			if !m.cmdStatus.IsStarted() {
			}
			if m.cmdStatus.IsFinished() {
				m.cmdStatus = m.cmdStatus.Reset()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return "Viewing the lights on/off screen - Esc to go back to main menu"
}
