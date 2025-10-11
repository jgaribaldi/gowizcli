package main

import (
	"flag"
	"gowizcli/client"
	"gowizcli/db"
	"gowizcli/luminance"
	"gowizcli/wiz"
)

func main() {
	var config Config
	readConfigFile(&config)
	readConfigEnvironment(&config)

	var command string
	var destAddress string

	flag.StringVar(&destAddress, "address", "255.255.255.255", "Destination address of the command - Use the local broadcast address for 'discover'")
	flag.StringVar(&command, "command", "", "Command to execute. Valid values are discover, show, reset, on, off")
	flag.Parse()

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

	cmd, err := client.NewCommand(command)
	if err != nil {
		panic(err)
	}
	cmd.AddParameters([]string{destAddress})

	c.Execute(*cmd)
}
