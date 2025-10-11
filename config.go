package main

import (
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
	Database struct {
		File string `yaml:"file"`
	} `yaml:"database"`
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
