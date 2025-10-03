package main

import (
	"fmt"
	"strings"
)

type CommandType int

const (
	Discover CommandType = iota
	Show
	Reset
	TurnOn
	TurnOff
)

var commandName = map[CommandType]string{
	Discover: "discover",
	Show:     "show",
	Reset:    "reset",
	TurnOn:   "on",
	TurnOff:  "off",
}

var commandMap = map[string]CommandType{
	"discover": Discover,
	"show":     Show,
	"reset":    Reset,
	"on":       TurnOn,
	"off":      TurnOff,
}

func (c CommandType) String() string {
	return commandName[c]
}

func ParseString(str string) (CommandType, bool) {
	c, ok := commandMap[strings.ToLower(str)]
	return c, ok
}

type Command struct {
	CommandType CommandType
	Parameters  []string
}

type Client struct {
	wiz *Wiz
	db  *DBConnection
}

func NewClient(timeoutSecs int) (*Client, error) {
	conn, err := NewConnection(timeoutSecs)
	if err != nil {
		return nil, err
	}

	db, err := NewDbConnection("lights.db")
	if err != nil {
		return nil, err
	}

	wiz := NewWiz(conn.Query)
	return &Client{
		wiz: wiz,
		db:  db,
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
		c.executeTurnOn()

	case TurnOff:
		c.executeTurnOff()

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
