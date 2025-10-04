package wiz

import (
	"encoding/json"
	"fmt"
	"gowizcli/infrastructure"

	"github.com/google/uuid"
)

type WizLight struct {
	Id         string
	MacAddress string
	IpAddress  string
}

type Wiz struct {
	query func(ipAddress string, message []byte) ([]infrastructure.QueryResponse, error)
}

func NewWiz(
	query func(ipAddress string, message []byte) ([]infrastructure.QueryResponse, error),
) *Wiz {
	return &Wiz{
		query: query,
	}
}

func (w Wiz) Discover(bcastAddr string) ([]WizLight, error) {
	fmt.Printf("Executing Wiz bulb discovery on network %s...\n", bcastAddr)

	getPilot := NewWizRequestBuilder().
		WithMethod("getPilot").
		Build()
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

	turnOn := NewWizRequestBuilder().
		WithMethod("setState").
		WithState(true).
		Build()
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

	turnOff := NewWizRequestBuilder().
		WithMethod("setState").
		WithState(false).
		Build()
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

type WizRequest struct {
	Id     int            `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
}

type WizRequestBuilder interface {
	WithMethod(method string) WizRequestBuilder
	WithDimming(dimming int) WizRequestBuilder
	WithRgb(r int, g int, b int) WizRequestBuilder
	WithTemp(temperature int) WizRequestBuilder
	WithSpeed(speed int) WizRequestBuilder
	WithScene(scene Scene) WizRequestBuilder
	WithState(state bool) WizRequestBuilder
	Build() *WizRequest
}

type wizRequestBuilder struct {
	wizRequest *WizRequest
}

func NewWizRequestBuilder() WizRequestBuilder {
	return &wizRequestBuilder{
		wizRequest: &WizRequest{
			Id:     1,
			Params: make(map[string]any),
		},
	}
}

func (w wizRequestBuilder) WithMethod(method string) WizRequestBuilder {
	w.wizRequest.Method = method
	return w
}

func (w wizRequestBuilder) WithDimming(dimming int) WizRequestBuilder {
	w.wizRequest.Params["dimming"] = dimming
	return w
}

func (w wizRequestBuilder) WithRgb(r int, g int, b int) WizRequestBuilder {
	w.wizRequest.Params["r"] = r
	w.wizRequest.Params["g"] = g
	w.wizRequest.Params["b"] = b
	return w
}

func (w wizRequestBuilder) WithTemp(temperature int) WizRequestBuilder {
	w.wizRequest.Params["temp"] = temperature
	return w
}

func (w wizRequestBuilder) WithSpeed(speed int) WizRequestBuilder {
	w.wizRequest.Params["speed"] = speed
	return w
}

func (w wizRequestBuilder) WithScene(scene Scene) WizRequestBuilder {
	w.wizRequest.Params["sceneId"] = scene
	return w
}

func (w wizRequestBuilder) WithState(state bool) WizRequestBuilder {
	w.wizRequest.Params["state"] = state
	return w
}

func (w wizRequestBuilder) Build() *WizRequest {
	return w.wizRequest
}

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
