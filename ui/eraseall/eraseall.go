package eraseall

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	client    *client.Client
	cmdStatus common.CmdStatus
	err       error
}

func NewModel(client *client.Client) Model {
	return Model{
		client:    client,
		cmdStatus: *common.NewCmdStatus(),
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			if m.cmdStatus.State == common.Running || m.cmdStatus.State == common.Done {
				return m, nil
			}
			m.cmdStatus = m.cmdStatus.Start()
			return m, eraseLightsCmd(m.client)
		}
	case eraseLightsMsg:
		m.err = nil
		m.cmdStatus = m.cmdStatus.Finish()
	case eraseLightsErrorMsg:
		m.err = msg.err
		m.cmdStatus = m.cmdStatus.Finish()
	}
	return m, nil
}

func (m Model) View() string {
	if m.cmdStatus.State == common.Running {
		return "Erasing all lights..."
	}

	if m.cmdStatus.State == common.Done {
		if m.err != nil {
			return fmt.Sprintf("Error erasing lights: %s - ESC to go back to main menu", m.err)
		}
		return "All lights erased - Execute a discover to populate the lights DB - ESC to go back to main menu"
	}

	return "Are you sure you want to delete all lights? y/N"
}

func eraseLightsCmd(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Reset,
			Parameters:  []string{},
		}

		_, err := c.Execute(cmd)
		if err != nil {
			return eraseLightsErrorMsg{err: err}
		}
		return eraseLightsMsg{}
	}
}

type eraseLightsMsg struct {
}

type eraseLightsErrorMsg struct {
	err error
}
