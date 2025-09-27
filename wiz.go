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
	connection *Connection
}

func NewWiz(conn *Connection) *Wiz {
	return &Wiz{
		connection: conn,
	}
}

func (w *Wiz) Discover() []WizLight {
	fmt.Printf("Executing Wiz bulb discovery on network %s...\n", w.connection.bcastAddr)

	getPilot := WizRequest{
		Method: "getPilot",
		Params: WizRequestParams{},
	}
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		println(err.Error())
		return []WizLight{}
	}

	srcIp, response, err := w.connection.Query(mGetPilot)
	if err != nil {
		fmt.Printf("Error executing query over the network: %s\n", err)
		return []WizLight{}
	}
	getPilotResult := WizResponse{}
	err = json.Unmarshal(response, &getPilotResult)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %s\n", err)
		return []WizLight{}
	}

	var result []WizLight
	result = make([]WizLight, 0)
	result = append(result, WizLight{
		MacAddress: getPilotResult.Result.Mac,
		IpAddress:  *srcIp,
	})
	return result
}

func (w *Wiz) TurnOn() {

}

func (w *Wiz) TurnOff() {

}
