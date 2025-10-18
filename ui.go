package main

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/wiz"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewType int

const (
	ViewMenu ViewType = iota
	ViewDiscover
	ViewShow
	ViewEraseAll
	ViewTurnOn
	ViewTurnOff
)

var viewCommandMap = map[client.CommandType]ViewType{
	client.Discover: ViewDiscover,
	client.Show:     ViewShow,
	client.Reset:    ViewEraseAll,
	client.TurnOn:   ViewTurnOn,
	client.TurnOff:  ViewTurnOff,
}

type model struct {
	currentView      ViewType
	viewHistory      []ViewType
	menuModel        MenuModel
	discoverModel    DiscoverModel
	showModel        ShowModel
	eraseAllModel    EraseAllModel
	lightsOnOffModel LightOnOffModel
	client           *client.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

type navigateToMsg struct {
	view ViewType
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if shouldQuit(msg) {
		return m, tea.Quit
	}

	if navMsg, ok := msg.(navigateToMsg); ok {
		m = navigateTo(m, navMsg.view)
		return m, nil
	}

	if shouldGoBack(msg) {
		m = navigateBack(m)
		return m, nil
	}

	return m.updateCurrentView(msg)
}

func (m model) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentView {
	case ViewMenu:
		m.menuModel, cmd = m.menuModel.Update(msg)
		return m, cmd
	case ViewDiscover:
		m.discoverModel, cmd = m.discoverModel.Update(msg)
		return m, cmd
	case ViewShow:
		m.showModel, cmd = m.showModel.Update(msg)
		return m, cmd
	case ViewEraseAll:
		m.eraseAllModel, cmd = m.eraseAllModel.Update(msg)
		return m, cmd
	case ViewTurnOn, ViewTurnOff:
		m.lightsOnOffModel, cmd = m.lightsOnOffModel.Update(msg)
		return m, cmd
	}
	return m, nil
}

func shouldQuit(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}

	keyName := key.String()
	return keyName == "ctrl+q"
}

func shouldGoBack(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}

	keyName := key.String()
	return keyName == "esc"
}

func (m model) View() string {
	switch m.currentView {
	case ViewMenu:
		return m.menuModel.View()
	case ViewDiscover:
		return m.discoverModel.View()
	case ViewShow:
		return m.showModel.View()
	case ViewEraseAll:
		return m.eraseAllModel.View()
	case ViewTurnOn, ViewTurnOff:
		return m.lightsOnOffModel.View()
	}
	return ""
}

func initialModel(client *client.Client) model {
	return model{
		currentView:      ViewMenu,
		viewHistory:      []ViewType{},
		menuModel:        NewMenuModel(),
		discoverModel:    NewDiscoverModel(client),
		showModel:        NewShowModel(client),
		eraseAllModel:    EraseAllModel{},
		lightsOnOffModel: LightOnOffModel{},
		client:           client,
	}
}

func navigateTo(m model, view ViewType) model {
	m.viewHistory = append(m.viewHistory, m.currentView)
	m.currentView = view
	return m
}

func navigateBack(m model) model {
	if len(m.viewHistory) > 0 {
		lastIndex := len(m.viewHistory) - 1
		m.currentView = m.viewHistory[lastIndex]
		m.viewHistory = m.viewHistory[:lastIndex]
	}
	return m
}

type DiscoverModel struct {
	client           *client.Client
	discovering      bool
	broadcastAddress string
	inputs           []textinput.Model
	focused          int
}

func NewDiscoverModel(client *client.Client) DiscoverModel {
	var inputs []textinput.Model
	inputs = make([]textinput.Model, 0)

	input1 := textinput.New()
	input1.Placeholder = "192"
	input1.CharLimit = 3
	input1.Width = 3
	input1.Prompt = ""
	input1.Validate = octetValidator
	input1.Focus()
	inputs = append(inputs, input1)

	input2 := textinput.New()
	input2.Placeholder = "168"
	input2.CharLimit = 3
	input2.Width = 3
	input2.Prompt = ""
	input2.Validate = octetValidator
	inputs = append(inputs, input2)

	input3 := textinput.New()
	input3.Placeholder = "1"
	input3.CharLimit = 3
	input3.Width = 3
	input3.Prompt = ""
	input3.Validate = octetValidator
	inputs = append(inputs, input3)

	input4 := textinput.New()
	input4.Placeholder = "255"
	input4.CharLimit = 3
	input4.Width = 3
	input4.Prompt = ""
	input4.Validate = octetValidator
	inputs = append(inputs, input4)

	return DiscoverModel{
		client:           client,
		discovering:      false,
		broadcastAddress: "",
		inputs:           inputs,
		focused:          0,
	}
}

func octetValidator(octet string) error {
	number, err := strconv.ParseInt(octet, 10, 64)
	if err != nil {
		return err
	}

	if number < 0 || number > 255 {
		return fmt.Errorf("incorrect octet: %s", octet)
	}

	return nil
}

func (m DiscoverModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DiscoverModel) Update(msg tea.Msg) (DiscoverModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m = nextInput(m)
			// return m, nil
		case tea.KeyShiftTab:
			m = previousInput(m)
			// return m, nil
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func previousInput(m DiscoverModel) DiscoverModel {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	return m
}

func nextInput(m DiscoverModel) DiscoverModel {
	m.focused = (m.focused + 1) % len(m.inputs)
	return m
}

func (m DiscoverModel) View() string {
	if m.broadcastAddress == "" {
		return fmt.Sprintf(
			"Broadcast address: %s.%s.%s.%s",
			m.inputs[0].View(),
			m.inputs[1].View(),
			m.inputs[2].View(),
			m.inputs[3].View(),
		)
	}
	if m.discovering {
		return fmt.Sprintf("Executing discovery on broadcast %s...", m.broadcastAddress)
	}
	return "Viewing the discovery screen - Esc to go back to main menu"
}

type EraseAllModel struct {
}

func (m EraseAllModel) Init() tea.Cmd {
	return nil
}

func (m EraseAllModel) Update(msg tea.Msg) (EraseAllModel, tea.Cmd) {
	return m, nil
}

func (m EraseAllModel) View() string {
	return "Viewing the erase all lights screen - Esc to go back to main menu"
}

type LightOnOffModel struct {
}

func (m LightOnOffModel) Init() tea.Cmd {
	return nil
}

func (m LightOnOffModel) Update(msg tea.Msg) (LightOnOffModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m LightOnOffModel) View() string {
	return "Viewing the lights on/off screen - Esc to go back to main menu"
}

type MenuModel struct {
	cursor   int
	options  []client.Option
	selected int
}

func NewMenuModel() MenuModel {
	var options []client.Option
	options = make([]client.Option, 0)

	for _, o := range client.Options {
		options = append(options, o)
	}

	return MenuModel{
		cursor:  0,
		options: options,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.selected = m.cursor
			fmt.Println("blah")
			selectedOption := m.options[m.selected]
			fmt.Printf("Selected option is %v\n", selectedOption)
			return m, func() tea.Msg {
				view := viewCommandMap[selectedOption.Type]
				return navigateToMsg{view: view}
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Welcome to gowizcli! A Wiz lights client written in Go. Select an option:\n\n")

	for i, o := range m.options {
		if m.cursor == i {
			s.WriteString("[*] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(o.Name)
		s.WriteString("\n")
	}

	return s.String()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ShowModel struct {
	table   table.Model
	loading bool
	err     error
	client  *client.Client
}

func NewShowModel(client *client.Client) ShowModel {
	columns := []table.Column{
		{Title: "ID", Width: 30},
		{Title: "MacAddress", Width: 20},
		{Title: "IpAddress", Width: 20},
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

type showDataLoadedMsg struct {
	lights []wiz.WizLight
}

type showDataErrorMsg struct {
	err error
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
		return "Fetching lights...\n"
	}
	return baseStyle.Render(m.table.View()) + "\n"
}
