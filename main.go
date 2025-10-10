package main

import (
	"flag"
	"gowizcli/client"
	"gowizcli/db"
	"gowizcli/wiz"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Luminance struct {
		IpGeolocation struct {
			ApiKey       string `yaml:"apiKey" envconfig:"IPGEOLOCATION_APIKEY"`
			Url          string `yaml:"url"`
			QueryTimeout int    `yaml:"queryTimeout"`
		} `yaml:"ipGeolocation"`
		OpenMeteo struct {
			Url          string `yaml:"url"`
			QueryTimeout int    `yaml:"queryTimeout"`
		} `yaml:"openMeteo"`
		Location struct {
			Latitude  float64 `yaml:"latitude"`
			Longitude float64 `yaml:"longitude"`
		} `yaml:"location"`
	} `yaml:"luminance"`
	Network struct {
		BroadcastAddress string `yaml:"broadcastAddress"`
		QueryTimeoutSec  int    `yaml:"queryTimeoutSec"`
	} `yaml:"network"`
}

func main() {
	var config Config
	readConfigFile(&config)
	readConfigEnvironment(&config)

	// ipGelocation := luminance.NewIpGeolocation(
	// 	config.Luminance.IpGeolocation.Url,
	// 	config.Luminance.IpGeolocation.ApiKey,
	// 	config.Luminance.IpGeolocation.QueryTimeout,
	// )
	// meteorology := luminance.NewMeteorology(
	// 	config.Luminance.OpenMeteo.Url,
	// 	config.Luminance.OpenMeteo.QueryTimeout,
	// )
	// orchestrator := luminance.NewOrchestrator(ipGelocation.GetSolarElevation, meteorology.GetCurrent)
	// luminance, err := orchestrator.GetCurrentLuminance(config.Luminance.Location.Latitude, config.Luminance.Location.Longitude)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Current luminance estimation is %f\n", luminance)

	var command string
	var destAddress string

	flag.StringVar(&destAddress, "address", "255.255.255.255", "Destination address of the command - Use the local broadcast address for 'discover'")
	flag.StringVar(&command, "command", "", "Command to execute. Valid values are discover, show, reset, on, off")
	flag.Parse()

	db, err := db.NewConnection("lights.db")
	if err != nil {
		panic(err)
	}

	conn, err := wiz.NewConnection(config.Network.QueryTimeoutSec)
	if err != nil {
		panic(err)
	}

	wiz := wiz.NewWiz(conn.Query)

	c, err := client.NewClient(
		db,
		wiz.Discover,
		wiz.TurnOn,
		wiz.TurnOff,
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

func readConfigFile(config *Config) error {
	file, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func readConfigEnvironment(config *Config) error {
	err := envconfig.Process("", config)
	if err != nil {
		return err
	}
	return nil
}
