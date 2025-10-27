package main

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/db"
	"gowizcli/luminance"
	"gowizcli/ui"
	"gowizcli/wiz"

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

	p := tea.NewProgram(ui.InitialModel(c), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v\n", err)
	}
}
