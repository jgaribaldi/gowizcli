package ui

import (
	"gowizcli/client"
	"gowizcli/wiz"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table      table.Model
	help       help.Model
	dimensions dimensions
	tableData  tableData
	cmdRunner  CmdRunner
}

func NewModel(client client.Functions) Model {
	columns := []table.Column{
		{Title: "IP Address", Width: 20},
		{Title: "MAC Address", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Tags", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
	)

	t.SetStyles(tableStyles())

	initialStatus := *NewCmdStatus()
	initialStatus = initialStatus.Start()

	return Model{
		table:     t,
		help:      help.New(),
		tableData: tableData{},
		cmdRunner: NewCmdRunner(client),
	}
}

func (m Model) Init() tea.Cmd {
	cmd := NewCmdRefresh(m.cmdRunner.client)
	_, t := m.cmdRunner.Run(cmd)
	return t
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case CmdDone:
		m.cmdRunner = m.cmdRunner.Finalize(msg)
		return m.handleCmdFinish(msg), nil
	case tea.KeyMsg:
		var cmd Command

		switch {
		case key.Matches(msg, keys.Refresh.binding):
			cmd = NewCmdRefresh(m.cmdRunner.client)
		case key.Matches(msg, keys.Switch.binding):
			if len(m.tableData.lights) == 0 {
				return m, nil
			}
			selected := m.table.Cursor()
			cmd = NewCmdSwitch(m.cmdRunner.client, m.tableData.lights[selected])
		case key.Matches(msg, keys.Discover.binding):
			cmd = NewCmdDiscover(m.cmdRunner.client)
		case key.Matches(msg, keys.EraseAll.binding):
			cmd = NewCmdEraseAll(m.cmdRunner.client)
		case key.Matches(msg, keys.Quit.binding):
			return m, tea.Quit
		}

		if cmd != nil {
			cr, t := m.cmdRunner.Run(cmd)
			m.cmdRunner = cr
			return m, t
		}
	case tea.WindowSizeMsg:
		return m.resize(msg), nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.cmdRunner.lastCmdStatus.State == Running {
		message := boxStyle.Render("Running command...")
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

func (m Model) handleCmdFinish(cmd CmdDone) Model {
	switch cmd.cmd.(type) {
	case CmdDiscover:
		m.tableData = tableData{
			err:    cmd.err,
			lights: merge(m.tableData.lights, cmd.lights),
		}
	case CmdSwitch:
		m.tableData = tableData{
			err:    cmd.err,
			lights: merge(m.tableData.lights, cmd.lights),
		}
	case CmdEraseAll:
		m.tableData = tableData{
			err:    cmd.err,
			lights: []wiz.Light{},
		}
	case CmdRefresh:
		m.tableData = tableData{
			err:    cmd.err,
			lights: merge(m.tableData.lights, cmd.lights),
		}
	}

	rows := make([]table.Row, len(m.tableData.lights))
	for i, l := range m.tableData.lights {
		rows[i] = lightToRow(l)
	}
	m.table.SetRows(rows)

	return m
}

func merge(existing []wiz.Light, incoming []wiz.Light) []wiz.Light {
	var existingIds = make(map[string]wiz.Light, len(existing))
	for _, l := range existing {
		existingIds[l.IpAddress] = l
	}

	var result = make([]wiz.Light, 0, len(existing)+len(incoming))
	seen := make(map[string]struct{}, len(incoming))

	for _, l := range incoming {
		result = append(result, l)
		seen[l.IpAddress] = struct{}{}
	}

	for _, l := range existing {
		if _, ok := seen[l.IpAddress]; !ok {
			result = append(result, l)
		}
	}
	return result
}

func lightToRow(l wiz.Light) table.Row {
	tagLine := strings.Join(l.Tags, ", ")
	if l.IsOn != nil {
		if *l.IsOn {
			return table.Row{
				l.IpAddress,
				parseMacAddress(l.MacAddress),
				"On",
				tagLine,
			}
		} else {
			return table.Row{
				l.IpAddress,
				parseMacAddress(l.MacAddress),
				"Off",
				tagLine,
			}
		}
	} else {
		return table.Row{
			l.IpAddress,
			parseMacAddress(l.MacAddress),
			"Unknown",
			tagLine,
		}
	}
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
