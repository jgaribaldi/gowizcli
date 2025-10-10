package client

import (
	"fmt"
	"gowizcli/db"
	"gowizcli/wiz"
)

type Client struct {
	lightsDb          db.LightsDatabase
	discoverAllLights func(string) ([]wiz.WizLight, error)
	turnOnLight       func(string) error
	turnOffLight      func(string) error
}

func NewClient(
	lightsDb db.LightsDatabase,
	discoverAllLights func(string) ([]wiz.WizLight, error),
	turnOnLight func(string) error,
	turnOffLight func(string) error,
) (*Client, error) {
	return &Client{
		lightsDb:          lightsDb,
		discoverAllLights: discoverAllLights,
		turnOnLight:       turnOnLight,
		turnOffLight:      turnOffLight,
	}, nil
}

func (c Client) Execute(command Command) error {
	switch command.CommandType {

	case Discover:
		c.executeDiscover(command.Parameters[0])

	case Show:
		c.executeShow()

	case Reset:
		c.executeReset()

	case TurnOn:
		c.executeTurnOn(command.Parameters[0])

	case TurnOff:
		c.executeTurnOff(command.Parameters[0])

	default:
		return fmt.Errorf("unknown command %s", command)
	}

	return nil
}

func (c Client) executeDiscover(bcastAddr string) error {
	lights, err := c.discoverAllLights(bcastAddr)
	if err != nil {
		return err
	}
	for _, light := range lights {
		fmt.Printf("Found new light with MAC Address %s and IP Address %s\n", light.MacAddress, light.IpAddress)
		_, err := c.lightsDb.Upsert(light)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) executeShow() error {
	fmt.Println("Lights")
	fmt.Println("------")

	lights, err := c.lightsDb.FindAll()
	if err != nil {
		return err
	}

	for _, l := range lights {
		fmt.Printf("ID: %s - MacAddress: %s - IpAddress: %s\n", l.Id, l.MacAddress, l.IpAddress)
	}
	return nil
}

func (c Client) executeReset() error {
	c.lightsDb.EraseAll()
	fmt.Println("Erased all data - run a discovery to populate again")
	return nil
}

func (c Client) executeTurnOn(destAddr string) error {
	return c.turnOnLight(destAddr)
}

func (c Client) executeTurnOff(destAddr string) error {
	return c.turnOffLight(destAddr)
}
