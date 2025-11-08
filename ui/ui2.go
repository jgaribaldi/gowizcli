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
		return m.resize(msg), nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

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

	cols := m.table.Columns()
	for i := range m.dimensions.columns {
		cols[i].Width = m.dimensions.columns[i].width
	}
	m.table.SetColumns(cols)

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
	window   size
	title    size
	table    size
	helpline size
	columns  []size
}

func (m Model) resize(msg tea.WindowSizeMsg) Model {
	windowSize := windowSize(msg)
	titleSize := titleSize(windowSize)
	helplineSize := helplineSize(windowSize, m)
	tableSize := tableSize(windowSize, titleSize, helplineSize, m)
	columnsSizes := columnsSize(tableSize, m)

	m.dimensions = dimensions{
		window:   windowSize,
		title:    titleSize,
		helpline: helplineSize,
		table:    tableSize,
		columns:  columnsSizes,
	}
	return m
}

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

func columnsSize(tableSize size, m Model) []size {
	columns := m.table.Columns()
	colNum := len(columns)
	if colNum == 0 {
		return nil
	}

	hPadding, _ := tableStyles().Cell.GetFrameSize()
	usableWidth := tableSize.width - (hPadding * colNum)
	usableWidth = max(usableWidth, colNum)

	base := usableWidth / colNum
	remainder := usableWidth % colNum

	var sizes []size = make([]size, colNum)
	for i := range columns {
		w := base
		// spread remainder to columns so total matches
		if i < remainder {
			w++
		}
		sizes[i] = size{
			width: w,
		}
	}

	return sizes
}

type size struct {
	width  int
	height int
}

var welcomeMsg string = "Welcome to Gowizcli! A Wiz client written in Go"
