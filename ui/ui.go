package ui

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"
	"strings"

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
	discoverStatus    common.CmdStatus
	eraseAllStatus    common.CmdStatus
	switchLightStatus common.CmdStatus
	table             table.Model
	help              help.Model
	dimensions        dimensions
	// lights            []wiz.Light
	tableData tableData
}

func NewModel(client *client.Client, bcastAddr string) Model {
	columns := []table.Column{
		{Title: "IP Address", Width: 20},
		{Title: "MAC Address", Width: 20},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
	)

	t.SetStyles(tableStyles())
	// var lights = make([]wiz.Light, 0)

	initialStatus := *common.NewCmdStatus()
	initialStatus = initialStatus.Start()

	return Model{
		client:            client,
		bcastAddr:         bcastAddr,
		fetchLigthsStatus: initialStatus,
		discoverStatus:    *common.NewCmdStatus(),
		eraseAllStatus:    *common.NewCmdStatus(),
		switchLightStatus: *common.NewCmdStatus(),
		table:             t,
		help:              help.New(),
		// lights:            lights,
		tableData: tableData{},
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchCmd()
}

func (m Model) update(lights []wiz.Light) Model {
	var newRows = make([]table.Row, len(lights))
	for idx, l := range lights {
		newRows[idx] = lightToRow(l)
	}

	m.table.SetRows(newRows)
	m.tableData.lights = lights
	m.table.Focus()

	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fetchDoneMsg:
		m.fetchLigthsStatus = m.fetchLigthsStatus.Finish()

		if msg.err != nil {
			m.tableData.err = msg.err
			return m, nil
		}
		m.tableData.err = nil
		return m.update(msg.lights), nil
	case switchDoneMsg:
		m.switchLightStatus = m.switchLightStatus.Finish()

		if msg.err != nil {
			m.tableData.err = msg.err
			return m, nil
		}

		m.tableData.err = nil
		var lights []wiz.Light = make([]wiz.Light, len(m.tableData.lights))
		copy(lights, m.tableData.lights)

		for idx, l := range m.tableData.lights {
			if l.Id == msg.light.Id {
				lights[idx] = msg.light
				break
			}
		}

		return m.update(lights), nil
	case discoverDoneMsg:
		m.discoverStatus = m.discoverStatus.Finish()
		if msg.err != nil || len(msg.lights) == 0 {
			m.tableData.err = msg.err
			return m, nil
		}
		m.tableData.err = nil
		return m.update(msg.lights), nil
	case eraseAllLightsDoneMsg:
		m.eraseAllStatus = m.eraseAllStatus.Finish()
		if msg.err != nil {
			m.tableData.err = msg.err
			return m, nil
		}
		m.tableData.err = nil
		return m.update([]wiz.Light{}), nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Refresh.binding):
			if m.fetchLigthsStatus.State != common.Running {
				m.fetchLigthsStatus = m.fetchLigthsStatus.Start()
				return m, m.fetchCmd()
			}
			return m, nil
		case key.Matches(msg, keys.Switch.binding):
			if m.switchLightStatus.State != common.Running {
				m.switchLightStatus = m.switchLightStatus.Start()
				return m, m.switchLightCmd()
			}
			return m, nil
		case key.Matches(msg, keys.Discover.binding):
			if m.discoverStatus.State != common.Running {
				m.discoverStatus = m.discoverStatus.Start()
				return m, m.discoverCommand()
			}
			return m, nil
		case key.Matches(msg, keys.EraseAll.binding):
			if m.eraseAllStatus.State != common.Running {
				m.eraseAllStatus = m.eraseAllStatus.Start()
				return m, m.eraseAllCommand()
			}
			return m, nil
		case key.Matches(msg, keys.Quit.binding):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m.resize(msg), nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func lightToRow(l wiz.Light) table.Row {
	if l.IsOn != nil {
		if *l.IsOn {
			return table.Row{
				l.IpAddress,
				parseMacAddress(l.MacAddress),
				"On",
			}
		} else {
			return table.Row{
				l.IpAddress,
				parseMacAddress(l.MacAddress),
				"Off",
			}
		}
	} else {
		return table.Row{
			l.IpAddress,
			parseMacAddress(l.MacAddress),
			"Unknown",
		}
	}
}

func (m Model) View() string {
	if m.fetchLigthsStatus.State == common.Running {
		message := boxStyle.Render("Fetching lights...")
		return lipgloss.Place(m.dimensions.window.width, m.dimensions.window.height, lipgloss.Center, lipgloss.Center, message)
	}

	if m.discoverStatus.State == common.Running {
		message := boxStyle.Render("Discovering lights...")
		return lipgloss.Place(m.dimensions.window.width, m.dimensions.window.height, lipgloss.Center, lipgloss.Center, message)
	}

	if m.eraseAllStatus.State == common.Running {
		message := boxStyle.Render("Erasing lights...")
		return lipgloss.Place(m.dimensions.window.width, m.dimensions.window.height, lipgloss.Center, lipgloss.Center, message)
	}

	if m.tableData.err != nil {
		// TODO: better handle error display
		return m.tableData.err.Error()
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

type keyAction struct {
	binding key.Binding
	run     func(*client.Client) (tea.Model, tea.Cmd)
}

type keyMap struct {
	Refresh  keyAction
	Switch   keyAction
	Discover keyAction
	EraseAll keyAction
	Quit     keyAction
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh.binding, k.Switch.binding, k.Discover.binding, k.EraseAll.binding, k.Quit.binding}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Refresh.binding, k.Switch.binding, k.Discover.binding, k.EraseAll.binding, k.Quit.binding},
	}
}

var keys = keyMap{
	Refresh:  keyAction{binding: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")), run: nil},
	Switch:   keyAction{binding: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Switch light")), run: nil},
	Discover: keyAction{binding: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Discover lights in network")), run: nil},
	EraseAll: keyAction{binding: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Erase all lights"))},
	Quit:     keyAction{binding: key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl+q", "Quit program")), run: nil},
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

func parseMacAddress(src string) string {
	positions := []int{1, 3, 5, 7, 9}

	var b strings.Builder
	next := 0
	for i, r := range src {
		b.WriteRune(r)

		if next < len(positions) && i == positions[next] {
			b.WriteByte(':')
			next++
		}
	}
	return b.String()
}

type tableData struct {
	lights []wiz.Light
	err    error
}

var welcomeMsg string = "Welcome to Gowizcli! A Wiz client written in Go"
