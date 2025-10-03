package main

import (
	"fmt"
	"strings"
)

type Command int

const (
	Discover Command = iota
	Show
	Reset
	TurnOn
	TurnOff
)

var commandName = map[Command]string{
	Discover: "discover",
	Show:     "show",
	Reset:    "reset",
	TurnOn:   "on",
	TurnOff:  "off",
}

var commandMap = map[string]Command{
	"discover": Discover,
	"show":     Show,
	"reset":    Reset,
	"on":       TurnOn,
	"off":      TurnOff,
}

func (c Command) String() string {
	return commandName[c]
}

func ParseString(str string) (Command, bool) {
	c, ok := commandMap[strings.ToLower(str)]
	return c, ok
}

type Client struct {
	wiz *Wiz
	db  *DBConnection
}

func NewClient(
	bcastAddr string,
	timeoutSecs int,
	query func(message []byte) ([]QueryResponse, error),
) (*Client, error) {
	conn, err := NewConnection(bcastAddr, timeoutSecs)
	if err != nil {
		return nil, err
	}

	db, err := NewDbConnection("lights.db")
	if err != nil {
		return nil, err
	}

	wiz := NewWiz(conn.Query, bcastAddr)
	return &Client{
		wiz: wiz,
		db:  db,
	}, nil
}

func (c Client) Execute(command Command) error {
	switch command {

	case Discover:
		c.executeDiscover()

	case Show:
		c.executeShow()

	case Reset:
		c.executeReset()

	case TurnOn:
		c.executeTurnOn()

	case TurnOff:
		c.executeTurnOff()

	default:
		return fmt.Errorf("unknown command %s", command)
	}

	return nil
}

func (c Client) executeDiscover() error {
	lights, err := c.wiz.Discover()
	if err != nil {
		return err
	}
	for _, light := range lights {
		fmt.Printf("Found new light with MAC Address %s and IP Address %s\n", light.MacAddress, light.IpAddress)
		_, err := c.db.Upsert(light)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) executeShow() error {
	fmt.Println("Lights")
	fmt.Println("------")

	lights, err := c.db.FindAll()
	if err != nil {
		return err
	}

	for _, l := range lights {
		fmt.Printf("ID: %s - MacAddress: %s - IpAddress: %s\n", l.Id, l.MacAddress, l.IpAddress)
	}
	return nil
}

func (c Client) executeReset() error {
	c.db.Reset()
	fmt.Println("Erased all data - run a discovery to populate again")
	return nil
}

func (c Client) executeTurnOn() error {
	return nil
}

func (c Client) executeTurnOff() error {
	return nil
}
