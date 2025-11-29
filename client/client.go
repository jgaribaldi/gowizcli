package client

import (
	"gowizcli/db"
	"gowizcli/luminance"
	"gowizcli/wiz"
)

type Location struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}

type Client struct {
	LightsDb  db.Storage
	WizClient wiz.Client
	Luminance luminance.Luminance
	Location  Location
}

type Functions interface {
	Discover() ([]wiz.Light, error)
	ShowAll() ([]wiz.Light, error)
	TurnOn(lightId string) (*wiz.Light, error)
	TurnOff(lightId string) (*wiz.Light, error)
	EraseAll()
}

func (c Client) Discover() ([]wiz.Light, error) {
	return c.executeDiscover()
}

func (c Client) ShowAll() ([]wiz.Light, error) {
	return c.executeShow()
}

func (c Client) TurnOn(lightId string) (*wiz.Light, error) {
	result, err := c.executeTurnOn(lightId)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (c Client) TurnOff(lightId string) (*wiz.Light, error) {
	result, err := c.executeTurnOff(lightId)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (c Client) EraseAll() {
	c.executeReset()
}

func (c Client) executeDiscover() ([]wiz.Light, error) {
	lights, err := c.WizClient.Discover()
	if err != nil {
		return nil, err
	}
	for _, light := range lights {
		_, err := c.LightsDb.Upsert(light)
		if err != nil {
			return nil, err
		}
	}
	return lights, nil
}

func (c Client) executeShow() ([]wiz.Light, error) {
	lights, err := c.LightsDb.FindAll()
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, len(lights))
	for i, l := range lights {
		result[i].Id = l.Id
		result[i].IpAddress = l.IpAddress
		result[i].MacAddress = l.MacAddress

		light, err := c.WizClient.Status(&l)
		if err == nil {
			result[i].IsOn = light.IsOn
		}
	}

	return result, nil
}

func (c Client) executeReset() ([]wiz.Light, error) {
	c.LightsDb.EraseAll()
	return nil, nil
}

func (c Client) executeTurnOn(lightId string) ([]wiz.Light, error) {
	light, err := c.LightsDb.FindById(lightId)
	if err != nil {
		return nil, err
	}

	newLight, err := c.WizClient.TurnOn(light)
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, 1)
	result[0] = *newLight
	return result, nil
}

func (c Client) executeTurnOff(lightId string) ([]wiz.Light, error) {
	light, err := c.LightsDb.FindById(lightId)
	if err != nil {
		return nil, err
	}

	newLight, err := c.WizClient.TurnOff(light)
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, 1)
	result[0] = *newLight
	return result, nil
}
