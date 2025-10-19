package showlights

import (
	"gowizcli/client"
	"gowizcli/wiz"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ShowModel struct {
	table   table.Model
	loading bool
	err     error
	client  *client.Client
}

func NewShowModel(client *client.Client) ShowModel {
	columns := []table.Column{
		{Title: "ID", Width: 40},
		{Title: "MAC Address", Width: 20},
		{Title: "IP Address", Width: 20},
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

	return ShowModel{
		table:   t,
		loading: true,
		err:     nil,
		client:  client,
	}
}

func (m ShowModel) Init() tea.Cmd {
	return fetchLightsCmd(m.client)
}

func (m ShowModel) Update(msg tea.Msg) (ShowModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case showDataLoadedMsg:
		rows := []table.Row{}
		for _, l := range msg.lights {
			rows = append(rows, table.Row{
				l.Id,
				l.MacAddress,
				l.IpAddress,
			})
		}
		m.table.SetRows(rows)
		m.loading = false
		m.err = nil
		return m, nil
	case showDataErrorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ShowModel) View() string {
	if m.loading {
		return "Fetching lights..."
	}

	if m.err != nil {
		return "Error fetching lights"
	}

	return baseStyle.Render(m.table.View()) + "\n"
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
	lights []wiz.WizLight
}

type showDataErrorMsg struct {
	err error
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
