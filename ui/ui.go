package ui

import (
	"gowizcli/client"
	"gowizcli/ui/discover"
	"gowizcli/ui/eraseall"
	"gowizcli/ui/showlights"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewType int

const (
	ViewMenu ViewType = iota
	ViewDiscover
	ViewShow
	ViewEraseAll
)

type model struct {
	currentView      ViewType
	viewHistory      []ViewType
	menuModel        MenuModel
	discoverModel    discover.Model
	showModel        showlights.Model
	eraseAllModel    eraseall.Model
	client           *client.Client
	defaultBcastAddr string
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
		return m.initCurrentView()
	}

	if shouldGoBack(msg) {
		m = navigateBack(m)
		return m, nil
	}

	return m.updateCurrentView(msg)
}

func (m model) initCurrentView() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewMenu:
		return m, m.menuModel.Init()
	case ViewDiscover:
		m.discoverModel = discover.NewModel(m.client, m.defaultBcastAddr)
		return m, m.discoverModel.Init()
	case ViewShow:
		m.showModel = showlights.NewModel(m.client)
		return m, m.showModel.Init()
	case ViewEraseAll:
		m.eraseAllModel = eraseall.NewModel(m.client)
		return m, m.eraseAllModel.Init()
	}

	return m, nil
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
	}
	return ""
}

func InitialModel(client *client.Client, defaultBcastAddr string) model {
	return model{
		currentView:      ViewMenu,
		viewHistory:      []ViewType{},
		menuModel:        NewMenuModel(),
		discoverModel:    discover.NewModel(client, defaultBcastAddr),
		showModel:        showlights.NewModel(client),
		eraseAllModel:    eraseall.NewModel(client),
		client:           client,
		defaultBcastAddr: defaultBcastAddr,
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
