package main

import (
	"gowizcli/client"
	"strings"

	"github.com/charmbracelet/bubbles/table"
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

type model struct {
	currentView ViewType
	viewHistory []ViewType

	menuModel        MenuModel
	discoverModel    DiscoverModel
	showModel        ShowModel2
	eraseAllModel    EraseAllModel
	lightsOnOffModel LightOnOffModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if shouldQuit(msg) {
		return m, tea.Quit
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

func initialModel() model {
	return model{
		currentView:      ViewMenu,
		viewHistory:      []ViewType{},
		menuModel:        MenuModel{},
		discoverModel:    DiscoverModel{},
		showModel:        ShowModel2{},
		eraseAllModel:    EraseAllModel{},
		lightsOnOffModel: LightOnOffModel{},
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
}

func (m DiscoverModel) Init() tea.Cmd {
	return nil
}

func (m DiscoverModel) Update(msg tea.Msg) (DiscoverModel, tea.Cmd) {
	return m, nil
}

func (m DiscoverModel) View() string {
	return "Viewing the discovery screen - Esc to go back to main menu"
}

type ShowModel2 struct {
}

func (m ShowModel2) Init() tea.Cmd {
	return nil
}

func (m ShowModel2) Update(msg tea.Msg) (ShowModel2, tea.Cmd) {
	return m, nil
}

func (m ShowModel2) View() string {
	return "Viewing the show lights screen - Esc to go back to main menu"
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
	return MenuModel{
		cursor:  0,
		options: client.Options,
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
		case "down", "j":
			m.cursor++
			if m.cursor >= len(client.Options) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(client.Options) - 1
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Welcome to gowizcli! A Wiz lights client written in Go. Pick a command:\n\n")
	for i, c := range client.Options {
		if m.cursor == i {
			s.WriteString("[*] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(c.Name)
		s.WriteString("\n")
	}
	return s.String()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ShowModel struct {
	table table.Model
}

func (m ShowModel) Init() tea.Cmd {
	return nil
}

func (m ShowModel) Update(msg tea.Msg) (ShowModel, tea.Cmd) {
	return m, nil
}

func (m ShowModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
