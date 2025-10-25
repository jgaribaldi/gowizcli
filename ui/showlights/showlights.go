package showlights

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	client    *client.Client
	data      showLightsData
	table     table.Model
	cmdStatus common.CommandStatus
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

	t.SetStyles(s)

	status := common.NewCommandStatus()
	status = status.Start()

	return Model{
		client:    client,
		data:      newShowLightsData(),
		table:     t,
		cmdStatus: status,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchLightsCmd(m.client)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case showDataLoadedMsg:
		m.cmdStatus = m.cmdStatus.Finish()
		m.data = m.data.result(msg.lights)

		var rows []table.Row = make([]table.Row, len(msg.lights))
		for idx, l := range msg.lights {
			rows[idx] = lightToRow(l)
		}
		m.table.SetRows(rows)

		return m, nil
	case showDataErrorMsg:
		m.cmdStatus = m.cmdStatus.Finish()
		m.data = m.data.error(msg.err)
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
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
	if m.cmdStatus.IsRunning() {
		return "Fetching lights..."
	}

	if m.data.err != nil {
		return "Error fetching lights"
	}

	return baseStyle.Render(m.table.View()) + "\n"
}

type showLightsData struct {
	lights []wiz.Light
	err    error
}

func newShowLightsData() showLightsData {
	return showLightsData{
		lights: make([]wiz.Light, 0),
		err:    nil,
	}
}

func (s showLightsData) result(lights []wiz.Light) showLightsData {
	for _, l := range lights {
		s.lights = append(s.lights, l)
	}
	s.err = nil
	return s
}

func (s showLightsData) error(err error) showLightsData {
	s.err = err
	s.lights = make([]wiz.Light, 0)
	return s
}

func fetchLightsCmd(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Show,
			Parameters:  []string{},
		}
		result, err := c.Execute(cmd)
		if err != nil {
			return showDataErrorMsg{err: err}
		}
		return showDataLoadedMsg{lights: result}
	}
}

type showDataLoadedMsg struct {
	lights []wiz.Light
}

type showDataErrorMsg struct {
	err error
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
