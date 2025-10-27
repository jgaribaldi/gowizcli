package wiz

import (
	"fmt"
	"testing"
)

func TestWizDiscover(t *testing.T) {
	var tests = []struct {
		response QueryResponse
		want     []Light
	}{
		{QueryResponse{
			SourceIpAddress: "192.168.1.174",
			Response:        []byte("{\"method\":\"getPilot\",\"env\":\"pro\",\"result\":{\"mac\":\"cc40857ce53c\",\"rssi\":-66,\"state\":true,\"sceneId\":8,\"speed\":100,\"dimming\":100}}"),
		}, []Light{
			{MacAddress: "cc40857ce53c", IpAddress: "192.168.1.174"},
		}},
	}

	for idx, tt := range tests {
		wiz := NewWiz(func(ipAddress string, message []byte) ([]QueryResponse, error) {
			var response []QueryResponse
			response = make([]QueryResponse, 0)
			response = append(response, tt.response)
			return response, nil
		})

		t.Run(fmt.Sprintf("Test %d", idx+1), func(t *testing.T) {
			got, _ := wiz.Discover("192.168.1.255")

			if got[0].IpAddress != tt.want[0].IpAddress || got[0].MacAddress != tt.want[0].MacAddress {
				t.Errorf("Got %s but want %s\n", got[0].IpAddress, tt.want[0].IpAddress)
			}
		})
	}
}
