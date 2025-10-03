package main

import (
	"flag"
	"fmt"
)

func main() {
	var command string
	var bcastAddr string
	var timeoutSecs int

	flag.StringVar(&bcastAddr, "broadcast", "255.255.255.255", "Broadcast address of the bulbs' local network")
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
		cmd.Parameters = []string{bcastAddr}
	}

	client.Execute(cmd)
}
