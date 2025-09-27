package main

import "flag"

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
		wiz.Discover()
	} else {
		println("Nothing to do")
	}

	println("Goodbye!")
}
