package main

import (
	"flag"
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

	cmd, err := client.NewCommand(command)
	if err != nil {
		panic(err)
	}
	cmd.AddParameters([]string{destAddress})

	c.Execute(*cmd)
}
