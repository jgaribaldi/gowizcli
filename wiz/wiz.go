package wiz

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Client interface {
	Discover(bcastAddr string) ([]Light, error)
	TurnOn(light *Light) (*Light, error)
	TurnOff(light *Light) (*Light, error)
	Status(light *Light) (*Light, error)
}

type Light struct {
	Id         string
	MacAddress string
	IpAddress  string
	IsOn       *bool
}

type Wiz struct {
	BulbClient  BulbClient
	TimeoutSecs int
}

func (w Wiz) Discover(bcastAddr string) ([]Light, error) {
	getPilot := NewRequestBuilder().
		WithMethod("getPilot").
		Build()
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: bcastAddr,
		Message:     mGetPilot,
		TimeoutSecs: w.TimeoutSecs,
	}
	queryResponse, err := w.BulbClient.Query(query)
	if err != nil {
		return nil, err
	}

	var result []Light = make([]Light, 0)
	for _, r := range queryResponse {
		getPilotResult := Response{}

		err = json.Unmarshal(r.Response, &getPilotResult)
		if err != nil {
			return nil, err
		} else {
			result = append(result, Light{
				Id:         uuid.New().String(),
				MacAddress: getPilotResult.Result.Mac,
				IpAddress:  r.Source,
			})
		}
	}

	return result, nil
}

func (w Wiz) TurnOn(light *Light) (*Light, error) {
	turnOn := NewRequestBuilder().
		WithMethod("setState").
		WithState(true).
		Build()
	mTurnOn, err := json.Marshal(turnOn)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: light.IpAddress,
		Message:     mTurnOn,
		TimeoutSecs: w.TimeoutSecs,
	}
	_, err = w.BulbClient.Query(query)
	if err != nil {
		return nil, err
	}

	return w.Status(light)
}

func (w Wiz) TurnOff(light *Light) (*Light, error) {
	turnOff := NewRequestBuilder().
		WithMethod("setState").
		WithState(false).
		Build()
	mTurnOff, err := json.Marshal(turnOff)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: light.IpAddress,
		Message:     mTurnOff,
		TimeoutSecs: w.TimeoutSecs,
	}
	_, err = w.BulbClient.Query(query)
	if err != nil {
		return nil, err
	}

	return w.Status(light)
}

func (w Wiz) Status(light *Light) (*Light, error) {
	getPilot := NewRequestBuilder().
		WithMethod("getPilot").
		Build()
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: light.IpAddress,
		Message:     mGetPilot,
		TimeoutSecs: w.TimeoutSecs,
	}
	queryResponse, err := w.BulbClient.Query(query)
	if err != nil {
		return nil, err
	}

	if len(queryResponse) > 0 {
		getPilotResult := Response{}
		err = json.Unmarshal(queryResponse[0].Response, &getPilotResult)
		if err != nil {
			return nil, err
		}

		return &Light{
			Id:         light.Id,
			MacAddress: light.MacAddress,
			IpAddress:  light.IpAddress,
			IsOn:       &getPilotResult.Result.State,
		}, nil
	}

	return nil, fmt.Errorf("device on address %s did not respond", light.IpAddress)
}
