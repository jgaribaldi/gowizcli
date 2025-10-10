package main

import (
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
	} `yaml:"luminance"`
}

func main() {
	var config Config
	readConfigFile(&config)
	readConfigEnvironment(&config)

	_ = luminance.NewIpGeolocation(
		config.Luminance.IpGeolocationUrl,
		config.Luminance.IpGeolocationApiKey,
		config.Luminance.IpGeolocationQueryTimeout,
	)
	_ = luminance.NewMeteorology(
		config.Luminance.OpenMeteoUrl,
		config.Luminance.OpenMeteoQueryTimeout,
	)

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
