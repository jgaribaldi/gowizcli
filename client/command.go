package client

import (
	"fmt"
	"strings"
)

type Option struct {
	Type CommandType
	Name string
}

var Options = []Option{
	{Discover, "Discover lights in local network"},
	{Show, "Show the discovered lights"},
	{Reset, "Delete all discovered lights"},
	{TurnOn, "Turn a light on"},
	{TurnOff, "Turn a light off"},
}

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

type Command struct {
	CommandType CommandType
	Parameters  []string
}

func NewCommand(cmdName string) (*Command, error) {
	c, ok := commandMap[strings.ToLower(cmdName)]
	if !ok {
		return nil, fmt.Errorf("unknown command %s", cmdName)
	}

	return &Command{
		CommandType: c,
		Parameters:  make([]string, 0),
	}, nil
}

func (c *Command) AddParameters(parameters []string) {
	switch c.CommandType {
	case Discover, TurnOn, TurnOff:
		c.Parameters = append(c.Parameters, parameters...)
	}
}
