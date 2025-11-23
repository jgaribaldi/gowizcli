package wiz

import (
	"fmt"
	"net"
	"time"
)

type NetworkConfig struct {
	BroadcastAddress string `yaml:"broadcastAddress"`
	QueryTimeoutSec  int    `yaml:"queryTimeoutSec"`
}

type BulbClient interface {
	Query(bulbQuery BulbQuery) ([]BulbResponse, error)
}

type BulbQuery struct {
	Destination string
	Message     []byte
	TimeoutSecs int
}

type BulbResponse struct {
	Source   string
	Response []byte
}

type UDPClient struct {
}

func (c UDPClient) Query(bulbQuery BulbQuery) ([]BulbResponse, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	address := fmt.Sprintf("%s:%s", bulbQuery.Destination, bulbPort)
	resolved, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	_, err = conn.WriteTo(bulbQuery.Message, resolved)
	if err != nil {
		return nil, err
	}

	var result []BulbResponse = make([]BulbResponse, 0)
	var buffer []byte = make([]byte, 1024)
	for {
		timeout := time.Now().Add(time.Duration(bulbQuery.TimeoutSecs) * time.Second)
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
			Source:   clientAddr.(*net.UDPAddr).IP.String(),
			Response: buffer[:n],
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
