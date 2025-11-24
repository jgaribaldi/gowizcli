package main

import (
	"gowizcli/client"
	"gowizcli/luminance"
	"gowizcli/wiz"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Luminance struct {
		IpGeolocation luminance.IpGeolocationConfig `yaml:"ipGeolocation"`
		OpenMeteo     luminance.OpenMeteoConfig     `yaml:"openMeteo"`
	} `yaml:"luminance"`
	Location client.Location   `yaml:"location"`
	Network  wiz.NetworkConfig `yaml:"network"`
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
