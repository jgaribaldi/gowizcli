package main

import (
	"gowizcli/client"
	"gowizcli/db"
	"gowizcli/luminance"
	"gowizcli/wiz"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var config Config
	readConfigFile(&config)
	readConfigEnvironment(&config)

	db, err := db.NewConnection(config.Database.File)
	if err != nil {
		panic(err)
	}

	conn, err := wiz.NewConnection(config.Network.QueryTimeoutSec)
	if err != nil {
		panic(err)
	}

	wiz := wiz.NewWiz(conn.Query)
	ipGelocation := luminance.NewIpGeolocation(
		config.Luminance.IpGeolocation.Url,
		config.Luminance.IpGeolocation.ApiKey,
		config.Luminance.IpGeolocation.QueryTimeout,
	)
	meteorology := luminance.NewMeteorology(
		config.Luminance.OpenMeteo.Url,
		config.Luminance.OpenMeteo.QueryTimeout,
	)
	orchestrator := luminance.NewOrchestrator(ipGelocation.GetSolarElevation, meteorology.GetCurrent)

	c, err := client.NewClient(
		db,
		wiz,
		orchestrator.GetCurrentLuminance,
	)
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(model{})
	m, err := p.Run()
	if err != nil {
		panic(err)
	}

	if m, ok := m.(model); ok && m.choice.Name != "" {
		cmd := client.Command{CommandType: m.choice.Type}
		c.Execute(cmd)
	}
}

type model struct {
	cursor int
	choice client.Option
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.choice = client.Options[m.cursor]
			return m, tea.Quit
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

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("Welcome to gowizcli! A Wiz lights client written in Go. Pick a command:\n\n")
	for i, c := range client.Options {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(c.Name)
		s.WriteString("\n")
	}
	return s.String()
}
