package wiz

import (
	"encoding/json"

	"github.com/google/uuid"
)

type WizClient interface {
	Discover(bcastAddr string) ([]WizLight, error)
	TurnOn(destAddr string) error
	TurnOff(destAddr string) error
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

	var result []WizLight
	result = make([]WizLight, 0)

	for _, r := range queryResponse {
		getPilotResult := Response{}

		err = json.Unmarshal(r.Response, &getPilotResult)
		if err != nil {
			return nil, err
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
	params := make(map[string]any)
	params["state"] = true

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
	params := make(map[string]any)
	params["state"] = false

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
