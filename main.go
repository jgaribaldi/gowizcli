package main

import (
	"flag"
	"fmt"
)

func main() {
	var discover bool
	var bcastAddr string
	var timeoutSecs int

	flag.BoolVar(&discover, "discover", false, "Execute a discovery on the given broadcast address")
	flag.StringVar(&bcastAddr, "broadcast", "255.255.255.255", "Broadcast address of the bulbs' local network")
	flag.IntVar(&timeoutSecs, "timeout", 1, "Query timeout in seconds")
	flag.Parse()

	conn, err := NewConnection(bcastAddr, timeoutSecs)
	if err != nil {
		panic(err)
	}
	wiz := NewWiz(conn)

	if discover {
		lights := wiz.Discover()
		for _, light := range lights {
			fmt.Printf("Found new light with MAC Address %s and IP Address %s\n", light.MacAddress, light.IpAddress)
		}
	} else {
		println("Nothing to do")
	}

	println("Goodbye!")
}
