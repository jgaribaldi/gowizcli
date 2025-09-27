package main

import (
	"fmt"
	"net"
	"time"
)

type Connection struct {
	bcastAddr        *net.UDPAddr
	queryTimeoutSecs int
}

func NewConnection(bcastAddr string, queryTimeoutSecs int) (*Connection, error) {
	address := fmt.Sprintf("%s:%s", bcastAddr, bulbPort)
	resolved, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	return &Connection{
		bcastAddr:        resolved,
		queryTimeoutSecs: queryTimeoutSecs,
	}, nil
}

func (c *Connection) Query(message string) (*string, []byte, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(c.queryTimeoutSecs) * time.Second))

	_, err = conn.WriteTo([]byte(message), c.bcastAddr)
	if err != nil {
		return nil, nil, err
	}

	buffer := make([]byte, 2048)
	n, clientAddr, err := conn.ReadFrom(buffer)
	if err != nil {
		return nil, nil, err
	}

	srcIpAddress := clientAddr.(*net.UDPAddr).IP.String()
	return &srcIpAddress, buffer[:n], nil
}

const bulbPort = "38899"
