package main

import (
	"flag"
	"fmt"
	"gowizcli/client"
)

func main() {
	var command string
	var destAddress string
	var timeoutSecs int

	flag.StringVar(&destAddress, "address", "255.255.255.255", "Destination address of the command - Use the local broadcast address for 'discover'")
	flag.IntVar(&timeoutSecs, "timeout", 1, "Query timeout in seconds")
	flag.StringVar(&command, "command", "", "Command to execute. Valid values are discover, show, reset, on, off")
	flag.Parse()

	c, err := client.NewClient(timeoutSecs)
	if err != nil {
		panic(err)
	}

	cmdType, ok := client.ParseString(command)
	if !ok {
		panic(fmt.Errorf("unknown command %s", command))
	}

	cmd := client.Command{}
	cmd.CommandType = cmdType
	switch cmdType {
	case client.Discover:
		cmd.Parameters = []string{destAddress}
	case client.TurnOff:
		cmd.Parameters = []string{destAddress}
	case client.TurnOn:
		cmd.Parameters = []string{destAddress}
	}

	c.Execute(cmd)
}
