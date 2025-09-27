package main

func main() {
	println("gowizcli")

	bcastAddr := "192.168.1.255"
	queryTimeoutSecs := 1
	conn, err := NewConnection(bcastAddr, queryTimeoutSecs)
	if err != nil {
		panic(err)
	}

	wiz := NewWiz(conn)
	wiz.Discover()
	println("that's all folks!")
}

