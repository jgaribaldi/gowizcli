package ui

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
	client            *client.Client
	bcastAddr         string
	fetchLigthsStatus common.CmdStatus
	fetchLightsData   fetchDoneMsg
	discoverData      discoverDoneMsg
	discoverStatus    common.CmdStatus
	table             table.Model
	help              help.Model
	dimensions        dimensions
}

func NewModel(client *client.Client, bcastAddr string) Model {
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
	)

	// s := table.DefaultStyles()
	// s.Header = s.Header.
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	BorderBottom(true).
	// 	Bold(false)
	// s.Selected = s.Selected.
	// 	Foreground(lipgloss.Color("229")).
	// 	Background(lipgloss.Color("57")).
	// 	Bold(true)

	t.SetStyles(tableStyles())

	return Model{
		client:            client,
		bcastAddr:         bcastAddr,
		fetchLigthsStatus: resetStatus(),
		fetchLightsData:   resetData(),
		discoverStatus:    *common.NewCmdStatus(),
		discoverData:      resetDiscoverData(),
		table:             t,
		help:              help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fetchDoneMsg:
		m.fetchLightsData = msg
		m.fetchLigthsStatus = m.fetchLigthsStatus.Finish()
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
			var rows []table.Row = make([]table.Row, len(m.fetchLightsData.lights))
			var lights []wiz.Light = make([]wiz.Light, len(m.fetchLightsData.lights))

			copy(rows, m.table.Rows())
			copy(lights, m.fetchLightsData.lights)

			for idx, l := range m.fetchLightsData.lights {
				if l.Id == msg.light.Id {
					rows[idx] = lightToRow(msg.light)
					lights[idx] = msg.light
					break
				}
			}

			m.table.SetRows(rows)
			m.fetchLightsData.lights = lights
			return m, nil
		}
	case discoverDoneMsg:
		if msg.err != nil || len(msg.lights) == 0 {
			return m, nil
		}
		m.discoverData = msg
		m.table.SetRows(rowsFromLights(m.discoverData.lights))
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Refresh):
			m.fetchLigthsStatus = resetStatus()
			m.fetchLightsData = resetData()
			return m, m.fetchCmd()
		case key.Matches(msg, keys.Switch):
			return m, m.switchLightCmd()
		case key.Matches(msg, keys.Discover):
			return m, m.discoverCommand()
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.resize(msg)
		// m.resizeTable()
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// func (m Model) resizeTable() {
// 	titleRendered := titleStyle.
// 		Width(m.dimensions.availableWidth).
// 		Render(welcomeMsg)
// 	titleHeight := lipgloss.Height(titleRendered)
//
// 	helpView := m.help.View(keys)
// 	helplineRendered := helplineStyle.
// 		Width(m.dimensions.availableWidth).
// 		Render(helpView)
// 	helplineHeight := lipgloss.Height(helplineRendered)
//
// 	m.table.SetHeight(0)
// 	tableRendered := tableStyle.
// 		Width(m.dimensions.availableWidth).
// 		Render(m.table.View())
// 	tableOverhead := lipgloss.Height(tableRendered)
//
// 	finalTableHeight := m.dimensions.availableHeight - titleHeight - helplineHeight - tableOverhead
// 	finalTableHeight = max(1, finalTableHeight)
//
// 	m.table.SetHeight(finalTableHeight)
// 	m.table.SetWidth(m.dimensions.availableWidth)
// }

func resetStatus() common.CmdStatus {
	status := common.NewCmdStatus()
	return status.Start()
}

func rowsFromLights(lights []wiz.Light) []table.Row {
	var rows []table.Row = make([]table.Row, 0)
	for _, l := range lights {
		rows = append(rows, lightToRow(l))
	}
	return rows
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
	if m.fetchLigthsStatus.State == common.Running {
		message := boxStyle.Render("Fetching lights...")
		return lipgloss.Place(m.dimensions.window.width, m.dimensions.window.height, lipgloss.Center, lipgloss.Center, message)
	}

	if m.fetchLightsData.err != nil {
		return "Error fetching lights"
	}

	if m.discoverStatus.State == common.Running {
		message := boxStyle.Render("Discovering lights...")
		return lipgloss.Place(m.dimensions.window.width, m.dimensions.window.height, lipgloss.Center, lipgloss.Center, message)
	}

	if m.discoverData.err != nil {
		// TODO: replace with modal
		return "Error discovering lights"
	}

	title := titleStyle.
		Width(m.dimensions.title.width).
		Render(welcomeMsg)

	helpView := m.help.View(keys)
	helpline := helplineStyle.
		Width(m.dimensions.helpline.width).
		Render(helpView)

	tableBody := tableStyle.
		Width(m.dimensions.table.width).
		Height(m.dimensions.table.height).
		Render(m.table.View())
	body := lipgloss.JoinVertical(lipgloss.Left, title, tableBody, helpline)
	return docStyle.Render(body)
}

type keyMap struct {
	Refresh  key.Binding
	Switch   key.Binding
	Discover key.Binding
	EraseAll key.Binding
	Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh, k.Switch, k.Discover, k.EraseAll, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh, k.Switch, k.Discover, k.EraseAll, k.Quit},
	}
}

var keys = keyMap{
	Refresh:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Switch:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Switch light")),
	Discover: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Discover lights in network")),
	EraseAll: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Erase all lights")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl+q", "Quit program")),
}

type dimensions struct {
	// totalHeight     int
	// totalWidth      int
	// availableHeight int
	// availableWidth  int
	// titleHeight     int
	// tableHeight     int
	// helplineHeight  int
	window   size
	title    size
	table    size
	helpline size
}

func (m Model) resize(msg tea.WindowSizeMsg) {
	// availableH, availableW := availableDimensions(msg.Height, msg.Width)
	windowSize := windowSize(msg)
	titleSize := titleSize(windowSize)
	helplineSize := helplineSize(windowSize, m)
	tableSize := tableSize(windowSize, titleSize, helplineSize, m)

	m.dimensions = dimensions{
		window:   windowSize,
		title:    titleSize,
		helpline: helplineSize,
		table:    tableSize,
		// totalHeight:     msg.Height,
		// totalWidth:      msg.Width,
		// availableHeight: availableH,
		// availableWidth:  availableW,
	}
}

// func availableDimensions(totalHeight, totalWidth int) (int, int) {
// 	wFrame, hFrame := docStyle.GetFrameSize()
//
// 	availableHeight := totalHeight - hFrame
// 	availableWidth := totalWidth - wFrame
//
// 	if availableHeight < 1 {
// 		availableHeight = 1
// 	}
//
// 	if availableWidth < 10 {
// 		availableWidth = 10
// 	}
//
// 	return availableHeight, availableWidth
// }

func windowSize(msg tea.WindowSizeMsg) size {
	marginWidth, marginHeight := docStyle.GetFrameSize()
	height := max(1, msg.Height-marginHeight)
	width := max(10, msg.Width-marginWidth)

	return size{
		width:  width,
		height: height,
	}
}

func titleSize(windowSize size) size {
	titleRendered := titleStyle.
		Width(windowSize.width).
		Render(welcomeMsg)
	titleHeight := lipgloss.Height(titleRendered)

	return size{
		width:  windowSize.width,
		height: titleHeight,
	}
}

func helplineSize(windowSize size, m Model) size {
	helpView := m.help.View(keys)
	helplineRendered := helplineStyle.
		Width(windowSize.width).
		Render(helpView)
	helplineHeight := lipgloss.Height(helplineRendered)

	return size{
		width:  windowSize.width,
		height: helplineHeight,
	}
}
func tableSize(windowSize, titleSize, helplineSize size, m Model) size {
	m.table.SetHeight(0)
	tableRendered := tableStyle.
		Width(windowSize.width).
		Render(m.table.View())
	tableOverhead := lipgloss.Height(tableRendered)

	finalTableHeight := windowSize.height - titleSize.height - helplineSize.height - tableOverhead
	finalTableHeight = max(1, finalTableHeight)

	return size{
		width:  windowSize.width,
		height: finalTableHeight,
	}
}

type size struct {
	width  int
	height int
}

var welcomeMsg string = "Welcome to Gowizcli! A Wiz client written in Go"
