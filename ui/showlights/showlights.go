package showlights

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	client    *client.Client
	cmdStatus common.CmdStatus
	data      fetchDoneMsg
	table     table.Model
	help      help.Model
}

func NewModel(client *client.Client) Model {
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
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	t.SetStyles(s)

	return Model{
		client:    client,
		cmdStatus: resetStatus(),
		data:      resetData(),
		table:     t,
		help:      help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchCmd()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fetchDoneMsg:
		m.data = msg
		m.cmdStatus = m.cmdStatus.Finish()
		m.table.Focus()

		if msg.err != nil {
			return m, nil
		} else {
			var rows []table.Row = make([]table.Row, len(msg.lights))
			for idx, l := range msg.lights {
				rows[idx] = lightToRow(l)
			}
			m.table.SetRows(rows)
			return m, nil
		}
	case switchDoneMsg:
		if msg.err != nil {
			return m, nil
		} else {
			var rows []table.Row = make([]table.Row, len(m.data.lights))
			var lights []wiz.Light = make([]wiz.Light, len(m.data.lights))

			copy(rows, m.table.Rows())
			copy(lights, m.data.lights)

			for idx, l := range m.data.lights {
				if l.Id == msg.light.Id {
					rows[idx] = lightToRow(msg.light)
					lights[idx] = msg.light
					break
				}
			}

			m.table.SetRows(rows)
			m.data.lights = lights
			return m, nil
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Refresh):
			m.cmdStatus = resetStatus()
			m.data = resetData()
			return m, m.fetchCmd()
		case key.Matches(msg, keys.Switch):
			return m, m.switchLightCmd()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func resetStatus() common.CmdStatus {
	status := common.NewCmdStatus()
	return status.Start()
}

func resetData() fetchDoneMsg {
	return fetchDoneMsg{
		lights: []wiz.Light{},
		err:    nil,
	}
}

func lightToRow(l wiz.Light) table.Row {
	if l.IsOn != nil {
		if *l.IsOn {
			return table.Row{
				l.Id,
				l.MacAddress,
				l.IpAddress,
				"On",
			}
		} else {
			return table.Row{
				l.Id,
				l.MacAddress,
				l.IpAddress,
				"Off",
			}
		}
	} else {
		return table.Row{
			l.Id,
			l.MacAddress,
			l.IpAddress,
			"Unknown",
		}
	}
}

func (m Model) View() string {
	if m.cmdStatus.State == common.Running {
		return "Fetching lights..."
	}

	if m.data.err != nil {
		return "Error fetching lights"
	}

	helpView := m.help.View(keys)
	return baseStyle.Render(m.table.View()) + "\n\n" + helpView
}

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
		if len(m.data.lights) > 0 {
			selectedRow := m.table.Cursor()
			if selectedRow < len(m.data.lights) {
				selectedLight := m.data.lights[selectedRow]

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
			return nil
		}
		return nil
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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type fetchDoneMsg struct {
	lights []wiz.Light
	err    error
}

type switchDoneMsg struct {
	light wiz.Light
	err   error
}

type keyMap struct {
	Refresh key.Binding
	Switch  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh, k.Switch}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Switch},
	}
}

var keys = keyMap{
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Switch:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Switch light")),
}
