package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWizDiscover(t *testing.T) {
	var tests = []struct {
		response QueryResponse
		want     []WizLight
	}{
		{QueryResponse{
			SourceIpAddress: "192.168.1.174",
			Response:        []byte("{\"method\":\"getPilot\",\"env\":\"pro\",\"result\":{\"mac\":\"cc40857ce53c\",\"rssi\":-66,\"state\":true,\"sceneId\":8,\"speed\":100,\"dimming\":100}}"),
		}, []WizLight{
			{MacAddress: "cc40857ce53c", IpAddress: "192.168.1.174"},
		}},
	}

	for idx, tt := range tests {
		wiz := NewWiz(func(message []byte) (*QueryResponse, error) {
			return &tt.response, nil
		}, "192.168.1.255")

		t.Run(fmt.Sprintf("Test %d", idx+1), func(t *testing.T) {
			got := wiz.Discover()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got %s but want %s\n", got, tt.want)
			}
		})
	}
}
