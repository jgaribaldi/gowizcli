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

	db, err := db.NewSQLiteDB(config.Database.File)
	if err != nil {
		panic(err)
	}

	bulbClient, err := wiz.NewUDPClient(config.Network.QueryTimeoutSec)
	if err != nil {
		panic(err)
	}

	wiz := wiz.NewWiz(bulbClient)
	astronomy := luminance.NewIpGeolocation(
		config.Luminance.IpGeolocation.Url,
		config.Luminance.IpGeolocation.ApiKey,
		config.Luminance.IpGeolocation.QueryTimeout,
	)
	meteorology := luminance.NewMeteorology(
		config.Luminance.OpenMeteo.Url,
		config.Luminance.OpenMeteo.QueryTimeout,
	)
	luminance := luminance.NewLuminance(astronomy, meteorology.GetCurrent)

	c, err := client.NewClient(
		db,
		wiz,
		config.Network.BroadcastAddress,
		luminance.GetCurrent,
	)
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(ui.NewModel(c), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v\n", err)
	}
}
