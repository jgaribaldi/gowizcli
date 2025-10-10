package main

import (
	"fmt"
	"gowizcli/luminance"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Luminance struct {
		IpGeolocationApiKey       string `yaml:"ipGeolocationApiKey" envconfig:"IPGEOLOCATION_APIKEY"`
		IpGeolocationUrl          string `yaml:"ipGeolocationUrl"`
		IpGeolocationQueryTimeout int    `yaml:"ipGeolocationQueryTimeout"`
		OpenMeteoUrl              string `yaml:"openMeteoUrl"`
		OpenMeteoQueryTimeout     int    `yaml:"openMeteoQueryTimeout"`
		Location                  struct {
			Latitude  float64 `yaml:"latitude"`
			Longitude float64 `yaml:"longitude"`
		} `yaml:"location"`
	} `yaml:"luminance"`
}

func main() {
	var config Config
	readConfigFile(&config)
	readConfigEnvironment(&config)

	ipGelocation := luminance.NewIpGeolocation(
		config.Luminance.IpGeolocationUrl,
		config.Luminance.IpGeolocationApiKey,
		config.Luminance.IpGeolocationQueryTimeout,
	)
	meteorology := luminance.NewMeteorology(
		config.Luminance.OpenMeteoUrl,
		config.Luminance.OpenMeteoQueryTimeout,
	)
	orchestrator := luminance.NewOrchestrator(ipGelocation.GetSolarElevation, meteorology.GetCurrent)
	luminance, err := orchestrator.GetCurrentLuminance(config.Luminance.Location.Latitude, config.Luminance.Location.Longitude)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current luminance estimation is %f\n", luminance)

	// var command string
	// var destAddress string
	// var timeoutSecs int

	// flag.StringVar(&destAddress, "address", "255.255.255.255", "Destination address of the command - Use the local broadcast address for 'discover'")
	// flag.IntVar(&timeoutSecs, "timeout", 1, "Query timeout in seconds")
	// flag.StringVar(&command, "command", "", "Command to execute. Valid values are discover, show, reset, on, off")
	// flag.Parse()

	// c, err := client.NewClient(timeoutSecs)
	// if err != nil {
	// 	panic(err)
	// }

	// cmd, err := client.NewCommand(command)
	// if err != nil {
	// 	panic(err)
	// }
	// cmd.AddParameters([]string{destAddress})

	// c.Execute(*cmd)
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
