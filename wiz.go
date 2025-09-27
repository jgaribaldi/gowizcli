package main

import "encoding/json"

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

type Wiz struct {
	connection *Connection
}

func NewWiz(conn *Connection) *Wiz {
	return &Wiz{
		connection: conn,
	}
}

func (w *Wiz) Discover() {
	getPilot := WizRequest{
		Method: "getPilot",
		Params: WizRequestParams{},
	}
	mGetPilot, err := json.Marshal(getPilot)
	if err != nil {
		println(err.Error())
		return
	}

	srcIp, response, err := w.connection.Query(mGetPilot)
	if err != nil {
		println(err.Error())
		return
	}
	getPilotResult := WizResponse{}
	err = json.Unmarshal(response, &getPilotResult)
	if err != nil {
		println(err.Error())
		return
	}
	println(*srcIp)
	println(getPilotResult.Result.Mac)
}

func (w *Wiz) TurnOn() {

}

func (w *Wiz) TurnOff() {

}
