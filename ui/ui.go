package ui

import (
	"fmt"
	"gowizcli/client"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

func InitialModel(client *client.Client) model {
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
		case tea.KeyShiftTab:
			m = previousInput(m)
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
