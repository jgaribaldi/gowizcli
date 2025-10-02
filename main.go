package main

import (
	"flag"
	"fmt"
)

func main() {
	var discover bool
	var bcastAddr string
	var timeoutSecs int
	var show bool
	var reset bool

	flag.BoolVar(&discover, "discover", false, "Execute a discovery on the given broadcast address")
	flag.StringVar(&bcastAddr, "broadcast", "255.255.255.255", "Broadcast address of the bulbs' local network")
	flag.IntVar(&timeoutSecs, "timeout", 1, "Query timeout in seconds")
	flag.BoolVar(&show, "show", false, "Show the contents of the lights database")
	flag.BoolVar(&reset, "reset", false, "Clear the lights database")
	flag.Parse()

	conn, err := NewConnection(bcastAddr, timeoutSecs)
	if err != nil {
		panic(err)
	}

	db, err := NewDbConnection("lights.db")
	if err != nil {
		panic(err)
	}

	wiz := NewWiz(conn.Query, bcastAddr)

	if discover {
		lights, err := wiz.Discover()
		if err != nil {
			panic(err)
		}
		for _, light := range lights {
			fmt.Printf("Found new light with MAC Address %s and IP Address %s\n", light.MacAddress, light.IpAddress)
			db.Upsert(light)
		}
	} else if show {
		println("Lights")
		println("------")

		lights, err := db.FindAll()
		if err != nil {
			panic(err)
		}

		for _, l := range lights {
			fmt.Printf("ID: %s - MacAddress: %s - IpAddress: %s\n", l.Id, l.MacAddress, l.IpAddress)
		}
	} else if reset {
		db.Reset()
		println("Erased all data - run a discovery to populate again")
	} else {
		println("Nothing to do")
	}

	println("Goodbye!")
}
