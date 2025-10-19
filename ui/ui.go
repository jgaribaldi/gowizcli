package ui

import (
	"gowizcli/client"

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
		return m, m.initCurrentView()
	}

	if shouldGoBack(msg) {
		m = navigateBack(m)
		return m, nil
	}

	return m.updateCurrentView(msg)
}

func (m model) initCurrentView() tea.Cmd {
	switch m.currentView {
	case ViewMenu:
		return m.menuModel.Init()
	case ViewDiscover:
		return m.discoverModel.Init()
	case ViewShow:
		return m.showModel.Init()
	case ViewEraseAll:
		return m.eraseAllModel.Init()
	case ViewTurnOn, ViewTurnOff:
		return m.lightsOnOffModel.Init()
	}

	return nil
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
