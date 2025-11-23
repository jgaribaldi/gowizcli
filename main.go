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

	wiz := wiz.Wiz{
		BulbClient:  wiz.UDPClient{},
		TimeoutSecs: config.Network.QueryTimeoutSec,
	}

	luminance := luminance.Luminance{
		Astronomy: luminance.IpGeolocation{
			Config: config.Luminance.IpGeolocation,
		},
		Meteorology: luminance.OpenMeteo{
			Config: config.Luminance.OpenMeteo,
		},
	}

	c := client.Client{
		LightsDb:  db,
		WizClient: wiz,
		Luminance: luminance,
		NetConfig: config.Network,
	}

	p := tea.NewProgram(ui.NewModel(&c), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v\n", err)
	}
}
