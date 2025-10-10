package client

import (
	"fmt"
	"gowizcli/wiz"
)

type Client struct {
	wiz            *wiz.Wiz
	upsertLight    func(wiz.WizLight) (*wiz.WizLight, error)
	findAllLights  func() ([]wiz.WizLight, error)
	eraseAllLights func()
}

func NewClient(
	timeoutSecs int,
	upsertLight func(wiz.WizLight) (*wiz.WizLight, error),
	findAllLights func() ([]wiz.WizLight, error),
	eraseAllLights func(),
) (*Client, error) {
	conn, err := wiz.NewConnection(timeoutSecs)
	if err != nil {
		return nil, err
	}

	wiz := wiz.NewWiz(conn.Query)
	return &Client{
		wiz:            wiz,
		upsertLight:    upsertLight,
		findAllLights:  findAllLights,
		eraseAllLights: eraseAllLights,
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
	lights, err := c.wiz.Discover(bcastAddr)
	if err != nil {
		return err
	}
	for _, light := range lights {
		fmt.Printf("Found new light with MAC Address %s and IP Address %s\n", light.MacAddress, light.IpAddress)
		_, err := c.upsertLight(light)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) executeShow() error {
	fmt.Println("Lights")
	fmt.Println("------")

	lights, err := c.findAllLights()
	if err != nil {
		return err
	}

	for _, l := range lights {
		fmt.Printf("ID: %s - MacAddress: %s - IpAddress: %s\n", l.Id, l.MacAddress, l.IpAddress)
	}
	return nil
}

func (c Client) executeReset() error {
	c.eraseAllLights()
	fmt.Println("Erased all data - run a discovery to populate again")
	return nil
}

func (c Client) executeTurnOn(destAddr string) error {
	return c.wiz.TurnOn(destAddr)
}

func (c Client) executeTurnOff(destAddr string) error {
	return c.wiz.TurnOff(destAddr)
}
