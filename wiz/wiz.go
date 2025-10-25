package wiz

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Client interface {
	Discover(bcastAddr string) ([]Light, error)
	TurnOn(destAddr string) error
	TurnOff(destAddr string) error
	IsTurnedOn(destAddr string) (*bool, error)
}

type Light struct {
	Id         string
	MacAddress string
	IpAddress  string
	IsOn       *bool
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

func (w Wiz) Discover(bcastAddr string) ([]Light, error) {
	getPilot := NewRequestBuilder().
		WithMethod("getPilot").
		Build()
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		return nil, err
	}

	queryResponse, err := w.query(bcastAddr, mGetPilot)
	if err != nil {
		return nil, err
	}

	var result []Light
	result = make([]Light, 0)

	for _, r := range queryResponse {
		getPilotResult := Response{}

		err = json.Unmarshal(r.Response, &getPilotResult)
		if err != nil {
			return nil, err
		} else {
			result = append(result, Light{
				Id:         uuid.New().String(),
				MacAddress: getPilotResult.Result.Mac,
				IpAddress:  r.SourceIpAddress,
			})
		}
	}

	return result, nil
}

func (w Wiz) TurnOn(destAddr string) error {
	turnOn := NewRequestBuilder().
		WithMethod("setState").
		WithState(true).
		Build()
	mTurnOn, err := json.Marshal(turnOn)
	if err != nil {
		return err
	}

	_, err = w.query(destAddr, mTurnOn)
	if err != nil {
		return err
	}

	return nil
}

func (w Wiz) TurnOff(destAddr string) error {
	turnOff := NewRequestBuilder().
		WithMethod("setState").
		WithState(false).
		Build()
	mTurnOff, err := json.Marshal(turnOff)
	if err != nil {
		return err
	}

	_, err = w.query(destAddr, mTurnOff)
	if err != nil {
		return err
	}

	return nil
}

func (w Wiz) IsTurnedOn(destAddr string) (*bool, error) {
	params := make(map[string]any)
	params["state"] = false

	isTurnedOn := NewRequestBuilder().
		WithMethod("getPilot").
		Build()
	mIsTurnedOn, err := json.Marshal(isTurnedOn)
	if err != nil {
		return nil, err
	}

	isTurnedOnResponse, err := w.query(destAddr, mIsTurnedOn)
	if err != nil {
		return nil, err
	}

	if len(isTurnedOnResponse) > 0 {
		response := Response{}
		err = json.Unmarshal(isTurnedOnResponse[0].Response, &response)
		if err != nil {
			return nil, err
		}

		isOn := response.Result.State
		return &isOn, nil
	} else {
		return nil, fmt.Errorf("no response from device")
	}
}
