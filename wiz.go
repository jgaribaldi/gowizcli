package main

import (
	"encoding/json"
	"fmt"
)

type WizRequestParams struct {
	Params map[string]string
}

type WizRequest struct {
	Method string           `json:"method"`
	Params WizRequestParams `json:"params"`
}

type WizResponseResult struct {
	Mac     string `json:"mac"`
	Rssi    int    `json:"rssi"`
	State   bool   `json:"state"`
	SceneId int    `json:"sceneId"`
	Temp    int    `json:"temp"`
	Dimming int    `json:"dimming"`
}

type WizResponse struct {
	Method string            `json:"method"`
	Env    string            `json:"env"`
	Result WizResponseResult `json:"result"`
}

type WizLight struct {
	MacAddress string
	IpAddress  string
}

type Wiz struct {
	query            func(message []byte) ([]QueryResponse, error)
	broadcastAddress string
}

func NewWiz(query func(message []byte) ([]QueryResponse, error), broadcastAddress string) *Wiz {
	return &Wiz{
		query:            query,
		broadcastAddress: broadcastAddress,
	}
}

func (w *Wiz) Discover() ([]WizLight, error) {
	fmt.Printf("Executing Wiz bulb discovery on network %s...\n", w.broadcastAddress)

	getPilot := WizRequest{
		Method: "getPilot",
		Params: WizRequestParams{},
	}
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		fmt.Printf("Error marshalling request: %s\n", err)
		return nil, err
	}

	queryResponse, err := w.query(mGetPilot)
	if err != nil {
		fmt.Printf("Error executing query over the network: %s\n", err)
		return nil, err
	}

	var result []WizLight
	result = make([]WizLight, 0)

	for _, r := range queryResponse {
		getPilotResult := WizResponse{}

		err = json.Unmarshal(r.Response, &getPilotResult)
		if err != nil {
			fmt.Printf("Error unmarshalling response: %s\n", err)
		} else {
			result = append(result, WizLight{
				MacAddress: getPilotResult.Result.Mac,
				IpAddress:  r.SourceIpAddress,
			})
		}
	}

	return result, nil
}

func (w *Wiz) TurnOn() {

}

func (w *Wiz) TurnOff() {

}
