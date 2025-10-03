package main

import (
	"fmt"
	"net"
	"time"
)

type QueryResponse struct {
	SourceIpAddress string
	Response        []byte
}

type Connection struct {
	queryTimeoutSecs int
}

func NewConnection(queryTimeoutSecs int) (*Connection, error) {
	return &Connection{
		queryTimeoutSecs: queryTimeoutSecs,
	}, nil
}

func (c *Connection) Query(ipAddress string, message []byte) ([]QueryResponse, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	address := fmt.Sprintf("%s:%s", ipAddress, bulbPort)
	resolved, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	_, err = conn.WriteTo(message, resolved)
	if err != nil {
		return nil, err
	}

	var result []QueryResponse
	result = make([]QueryResponse, 0)

	buffer := make([]byte, 1024)
	for {
		timeout := time.Now().Add(time.Duration(c.queryTimeoutSecs) * time.Second)
		err = conn.SetReadDeadline(timeout)
		if err != nil {
			return nil, err
		}

		n, clientAddr, err := conn.ReadFrom(buffer)
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			break
		}
		if err != nil {
			return nil, err
		}

		result = append(result, QueryResponse{
			SourceIpAddress: clientAddr.(*net.UDPAddr).IP.String(),
			Response:        buffer[:n],
		})
	}

	return result, nil
}

const bulbPort = "38899"
