package client

import (
	"fmt"
	"gowizcli/db"
	"gowizcli/wiz"
)

type Client struct {
	lightsDb     db.LightsDatabase
	wizClient    wiz.WizClient
	getLuminance func(float64, float64) (float64, error)
}

func NewClient(
	lightsDb db.LightsDatabase,
	wizClient wiz.WizClient,
	getLuminance func(float64, float64) (float64, error),
) (*Client, error) {
	return &Client{
		lightsDb:     lightsDb,
		wizClient:    wizClient,
		getLuminance: getLuminance,
	}, nil
}

func (c Client) Execute(command Command) ([]wiz.WizLight, error) {
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

func (c Client) executeDiscover(bcastAddr string) ([]wiz.WizLight, error) {
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

func (c Client) executeShow() ([]wiz.WizLight, error) {
	lights, err := c.lightsDb.FindAll()
	if err != nil {
		return nil, err
	}

	return lights, nil
}

func (c Client) executeReset() ([]wiz.WizLight, error) {
	c.lightsDb.EraseAll()
	return nil, nil
}

func (c Client) executeTurnOn(destAddr string) ([]wiz.WizLight, error) {
	return nil, c.wizClient.TurnOn(destAddr)
}

func (c Client) executeTurnOff(destAddr string) ([]wiz.WizLight, error) {
	return nil, c.wizClient.TurnOff(destAddr)
}
