package client

import "strings"

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
