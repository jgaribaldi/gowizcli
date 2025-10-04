package main

import (
	"flag"
	"fmt"
)

func main() {
	var command string
	var destAddress string
	var timeoutSecs int

	flag.StringVar(&destAddress, "address", "255.255.255.255", "Destination address of the command - Use the local broadcast address for 'discover'")
	flag.IntVar(&timeoutSecs, "timeout", 1, "Query timeout in seconds")
	flag.StringVar(&command, "command", "", "Command to execute. Valid values are discover, show, reset, on, off")
	flag.Parse()

	client, err := NewClient(timeoutSecs)
	if err != nil {
		panic(err)
	}

	cmdType, ok := ParseString(command)
	if !ok {
		panic(fmt.Errorf("unknown command %s", command))
	}

	cmd := Command{}
	cmd.CommandType = cmdType
	switch cmdType {
	case Discover:
		cmd.Parameters = []string{destAddress}
	case TurnOff:
		cmd.Parameters = []string{destAddress}
	case TurnOn:
		cmd.Parameters = []string{destAddress}
	}

	client.Execute(cmd)
}
