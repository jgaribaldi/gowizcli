package main

import "flag"

func main() {
	var discover bool
	flag.BoolVar(&discover, "discover", false, "Use this argument to execute a discovery on your local network")
	flag.Parse()

	if discover {
		bcastAddr := "192.168.1.255"
		queryTimeoutSecs := 1
		conn, err := NewConnection(bcastAddr, queryTimeoutSecs)
		if err != nil {
			panic(err)
		}

		wiz := NewWiz(conn)
		wiz.Discover()
	} else {
		println("Nothing to do")
	}

	println("Goodbye!")
}
