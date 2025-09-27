package main

func main() {
	println("gowizcli")

	wiz, err := NewWiz("192.168.1.255", 1)
	if err != nil {
		panic(err)
	}
	wiz.Discover()
	println("that's all folks!")
}

