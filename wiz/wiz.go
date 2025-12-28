package wiz

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Scene int

const (
	Ocean        Scene = 1
	Romance      Scene = 2
	Sunset       Scene = 3
	Party        Scene = 4
	Fireplace    Scene = 5
	Cozy         Scene = 6
	Forest       Scene = 7
	PastelColors Scene = 8
	WakeUp       Scene = 9
	Bedtime      Scene = 10
	WarmWhite    Scene = 11
	Daylight     Scene = 12
	CoolWhite    Scene = 13
	NightLight   Scene = 14
	Focus        Scene = 15
	Relax        Scene = 16
	TrueColors   Scene = 17
	TVTime       Scene = 18
	PlantGrowth  Scene = 19
	Spring       Scene = 20
	Summer       Scene = 21
	Fall         Scene = 22
	DeepDive     Scene = 23
	Jungle       Scene = 24
	Mojito       Scene = 25
	Club         Scene = 26
	Christmas    Scene = 27
	Halloween    Scene = 28
	Candlelight  Scene = 29
	GoldenWhite  Scene = 30
	Pulse        Scene = 31
	Steampunk    Scene = 32
	Rhythm       Scene = 1000
)

type Client interface {
	Discover() ([]Light, error)
	TurnOn(light *Light) (*Light, error)
	TurnOff(light *Light) (*Light, error)
	Status(light *Light) (*Light, error)
}

type Light struct {
	Id         string
	MacAddress string
	IpAddress  string
	IsOn       *bool
	Tags       []string
}

type Wiz struct {
	BulbClient BulbClient
	NetConfig  NetworkConfig
}

func (w Wiz) Discover() ([]Light, error) {
	getPilot := NewRequestBuilder().
		WithMethod("getPilot").
		Build()
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: w.NetConfig.BroadcastAddress,
		Message:     mGetPilot,
		TimeoutSecs: w.NetConfig.QueryTimeoutSec,
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
		TimeoutSecs: w.NetConfig.QueryTimeoutSec,
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
		TimeoutSecs: w.NetConfig.QueryTimeoutSec,
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
		TimeoutSecs: w.NetConfig.QueryTimeoutSec,
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

func (w Wiz) SetScene(light *Light, scene Scene) (*Light, error) {
	setScene := NewRequestBuilder().
		WithMethod("setPilot").
		WithScene(scene).
		Build()
	mSetScene, err := json.Marshal(setScene)
	if err != nil {
		return nil, err
	}

	query := BulbQuery{
		Destination: light.IpAddress,
		Message:     mSetScene,
		TimeoutSecs: w.NetConfig.QueryTimeoutSec,
	}
	_, err = w.BulbClient.Query(query)
	if err != nil {
		return nil, err
	}

	return w.Status(light)
}
