package wiz

import (
	"fmt"
	"net"
	"time"
)

type BulbClient interface {
	Query(ipAddress string, message []byte) ([]BulbResponse, error)
}

type BulbResponse struct {
	SourceIpAddress string
	Response        []byte
}

type UDPClient struct {
	queryTimeoutSecs int
}

func NewUDPClient(queryTimeoutSecs int) (*UDPClient, error) {
	return &UDPClient{
		queryTimeoutSecs: queryTimeoutSecs,
	}, nil
}

func (c UDPClient) Query(ipAddress string, message []byte) ([]BulbResponse, error) {
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

	var result []BulbResponse = make([]BulbResponse, 0)
	var buffer []byte = make([]byte, 1024)
	for {
		timeout := time.Now().Add(time.Duration(c.queryTimeoutSecs) * time.Second)
		err = conn.SetReadDeadline(timeout)
		if err != nil {
			return nil, err
		}

		n, clientAddr, err := conn.ReadFrom(buffer)
		if gotTimeout(err) {
			break
		}
		if err != nil {
			return nil, err
		}

		result = append(result, BulbResponse{
			SourceIpAddress: clientAddr.(*net.UDPAddr).IP.String(),
			Response:        buffer[:n],
		})
	}

	return result, nil
}

func gotTimeout(err error) bool {
	if err != nil {
		ne, ok := err.(net.Error)
		return ok && ne.Timeout()
	}
	return false
}

const bulbPort = "38899"
