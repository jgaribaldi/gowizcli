package client

import (
	"fmt"
	"gowizcli/db"
	"gowizcli/luminance"
	"gowizcli/wiz"
	"strings"
)

type Client struct {
	lightsDb  db.Storage
	wizClient wiz.Client
	bcastAddr string
	luminance luminance.Luminance
}

func NewClient(lightsDb db.Storage, wizClient wiz.Client, bcastAddr string, luminance luminance.Luminance) Client {
	return Client{
		lightsDb:  lightsDb,
		wizClient: wizClient,
		bcastAddr: bcastAddr,
		luminance: luminance,
	}
}

func (c Client) Execute(command Command) ([]wiz.Light, error) {
	switch command.CommandType {

	case Discover:
		return c.executeDiscover()

	case Show:
		return c.executeShow()

	case Reset:
		return c.executeReset()

	case TurnOn:
		return c.executeTurnOn(command.Parameters[0])

	case TurnOff:
		return c.executeTurnOff(command.Parameters[0])

	default:
		return nil, fmt.Errorf("unknown command %s", command)
	}
}

type CommandType int

const (
	Discover CommandType = iota
	Show
	Reset
	TurnOn
	TurnOff
)

type Command struct {
	CommandType CommandType
	Parameters  []string
}

func (c Command) String() string {
	params := strings.Join(c.Parameters, ", ")
	switch c.CommandType {
	case Discover:
		return "Discover " + params
	case Show:
		return "Show " + params
	case Reset:
		return "Reset " + params
	case TurnOn:
		return "Turn On " + params
	case TurnOff:
		return "Turn Off " + params
	}
	return ""
}

func (c Client) executeDiscover() ([]wiz.Light, error) {
	lights, err := c.wizClient.Discover(c.bcastAddr)
	if err != nil {
		return nil, err
	}
	for _, light := range lights {
		_, err := c.lightsDb.Upsert(light)
		if err != nil {
			return nil, err
		}
	}
	return lights, nil
}

func (c Client) executeShow() ([]wiz.Light, error) {
	lights, err := c.lightsDb.FindAll()
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, len(lights))
	for i, l := range lights {
		result[i].Id = l.Id
		result[i].IpAddress = l.IpAddress
		result[i].MacAddress = l.MacAddress

		light, err := c.wizClient.Status(&l)
		if err == nil {
			result[i].IsOn = light.IsOn
		}
	}

	return result, nil
}

func (c Client) executeReset() ([]wiz.Light, error) {
	c.lightsDb.EraseAll()
	return nil, nil
}

func (c Client) executeTurnOn(lightId string) ([]wiz.Light, error) {
	light, err := c.lightsDb.FindById(lightId)
	if err != nil {
		return nil, err
	}

	newLight, err := c.wizClient.TurnOn(light)
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, 1)
	result[0] = *newLight
	return result, nil
}

func (c Client) executeTurnOff(lightId string) ([]wiz.Light, error) {
	light, err := c.lightsDb.FindById(lightId)
	if err != nil {
		return nil, err
	}

	newLight, err := c.wizClient.TurnOff(light)
	if err != nil {
		return nil, err
	}

	var result []wiz.Light = make([]wiz.Light, 1)
	result[0] = *newLight
	return result, nil
}
