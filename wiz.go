package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type WizRequest struct {
	Id     int            `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
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
	Id         string
	MacAddress string
	IpAddress  string
}

type Wiz struct {
	query func(ipAddress string, message []byte) ([]QueryResponse, error)
}

func NewWiz(
	query func(ipAddress string, message []byte) ([]QueryResponse, error),
) *Wiz {
	return &Wiz{
		query: query,
	}
}

func (w Wiz) Discover(bcastAddr string) ([]WizLight, error) {
	fmt.Printf("Executing Wiz bulb discovery on network %s...\n", bcastAddr)

	getPilot := WizRequest{
		Id:     1,
		Method: "getPilot",
		Params: map[string]any{},
	}
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		fmt.Printf("Error marshalling request: %s\n", err)
		return nil, err
	}

	queryResponse, err := w.query(bcastAddr, mGetPilot)
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
				Id:         uuid.New().String(),
				MacAddress: getPilotResult.Result.Mac,
				IpAddress:  r.SourceIpAddress,
			})
		}
	}

	return result, nil
}

func (w Wiz) TurnOn(destAddr string) error {
	fmt.Printf("Turning on bulb with IP %s...\n", destAddr)

	params := make(map[string]any)
	params["state"] = true

	turnOn := WizRequest{
		Id:     1,
		Method: "setState",
		Params: params,
	}
	mTurnOn, err := json.Marshal(turnOn)
	if err != nil {
		fmt.Printf("Error marshalling request: %s\n", err)
		return err
	}

	_, err = w.query(destAddr, mTurnOn)
	if err != nil {
		fmt.Printf("Error executing query over the network: %s\n", err)
		return err
	}

	return nil
}

func (w Wiz) TurnOff(destAddr string) error {
	fmt.Printf("Turning off bulb with IP %s...\n", destAddr)
	params := make(map[string]any)
	params["state"] = false

	turnOff := WizRequest{
		Id:     1,
		Method: "setState",
		Params: params,
	}
	mTurnOff, err := json.Marshal(turnOff)
	if err != nil {
		fmt.Printf("Error marshalling request: %s\n", err)
		return err
	}

	_, err = w.query(destAddr, mTurnOff)
	if err != nil {
		fmt.Printf("Error executing query over the network: %s\n", err)
		return err
	}

	return nil
}
