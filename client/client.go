package client

import (
	"fmt"
	"gowizcli/db"
	"gowizcli/wiz"
)

type Client struct {
	lightsDb     db.LightsDatabase
	wizClient    wiz.Client
	getLuminance func(float64, float64) (float64, error)
}

func NewClient(
	lightsDb db.LightsDatabase,
	wizClient wiz.Client,
	getLuminance func(float64, float64) (float64, error),
) (*Client, error) {
	return &Client{
		lightsDb:     lightsDb,
		wizClient:    wizClient,
		getLuminance: getLuminance,
	}, nil
}

func (c Client) Execute(command Command) ([]wiz.Light, error) {
	switch command.CommandType {

	case Discover:
		return c.executeDiscover(command.Parameters[0])

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

func (c Client) executeDiscover(bcastAddr string) ([]wiz.Light, error) {
	lights, err := c.wizClient.Discover(bcastAddr)
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
	for idx, l := range lights {
		result[idx].Id = l.Id
		result[idx].IpAddress = l.IpAddress
		result[idx].MacAddress = l.MacAddress

		isOn, err := c.wizClient.IsTurnedOn(l.IpAddress)
		if err == nil {
			result[idx].IsOn = isOn
		}
	}

	return result, nil
}

func (c Client) executeReset() ([]wiz.Light, error) {
	c.lightsDb.EraseAll()
	return nil, nil
}

func (c Client) executeTurnOn(destAddr string) ([]wiz.Light, error) {
	return nil, c.wizClient.TurnOn(destAddr)
}

func (c Client) executeTurnOff(destAddr string) ([]wiz.Light, error) {
	return nil, c.wizClient.TurnOff(destAddr)
}
