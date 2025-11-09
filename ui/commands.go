package ui

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/wiz"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) fetchCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Show,
			Parameters:  []string{},
		}
		result, err := m.client.Execute(cmd)
		return fetchDoneMsg{
			lights: result,
			err:    err,
		}
	}
}

func (m Model) switchLightCmd() tea.Cmd {
	return func() tea.Msg {
		if len(m.tableData.lights) > 0 {
			selectedRow := m.table.Cursor()
			if selectedRow < len(m.tableData.lights) {
				selectedLight := m.tableData.lights[selectedRow]

				cmd := switchCommand(selectedLight)
				result, err := m.client.Execute(cmd)

				if len(result) > 0 {
					return switchDoneMsg{
						light: result[0],
						err:   err,
					}
				} else {
					return switchDoneMsg{
						err: err,
					}
				}
			}
			return switchDoneMsg{
				err: fmt.Errorf("invalid selected row"),
			}
		}
		return switchDoneMsg{
			err: fmt.Errorf("no lights to turn off/on"),
		}
	}
}

func switchCommand(light wiz.Light) client.Command {
	if light.IsOn != nil && *light.IsOn {
		return client.Command{
			CommandType: client.TurnOff,
			Parameters: []string{
				light.Id,
			},
		}

	} else {
		return client.Command{
			CommandType: client.TurnOn,
			Parameters: []string{
				light.Id,
			},
		}
	}
}

func (m Model) discoverCommand() tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Discover,
			Parameters: []string{
				m.bcastAddr,
			},
		}
		result, err := m.client.Execute(cmd)
		return discoverDoneMsg{
			lights: result,
			err:    err,
		}
	}
}

func (m Model) eraseAllCommand() tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Reset,
			Parameters:  []string{},
		}

		_, err := m.client.Execute(cmd)
		return eraseAllLightsDoneMsg{
			err: err,
		}
	}
}
