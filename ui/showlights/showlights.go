package showlights

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"

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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Refresh):
			m.cmdStatus = resetStatus()
			m.data = resetData()
			return m, m.fetchCmd()
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

	return baseStyle.Render(m.table.View()) + "\n"
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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type fetchDoneMsg struct {
	lights []wiz.Light
	err    error
}

var keyMap = struct {
	Refresh key.Binding
}{
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
}
