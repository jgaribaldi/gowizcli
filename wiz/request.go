package wiz

type Request struct {
	Id     int            `json:"id"`
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
}

type RequestBuilder interface {
	WithMethod(method string) RequestBuilder
	WithDimming(dimming int) RequestBuilder
	WithRgb(r int, g int, b int) RequestBuilder
	WithTemp(temperature int) RequestBuilder
	WithSpeed(speed int) RequestBuilder
	WithScene(scene Scene) RequestBuilder
	WithState(state bool) RequestBuilder
	Build() *Request
}

type requestBuilder struct {
	request *Request
}

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{
		request: &Request{
			Id:     1,
			Params: make(map[string]any),
		},
	}
}

func (w requestBuilder) WithMethod(method string) RequestBuilder {
	w.request.Method = method
	return w
}

func (w requestBuilder) WithDimming(dimming int) RequestBuilder {
	w.request.Params["dimming"] = dimming
	return w
}

func (w requestBuilder) WithRgb(r int, g int, b int) RequestBuilder {
	w.request.Params["r"] = r
	w.request.Params["g"] = g
	w.request.Params["b"] = b
	return w
}

func (w requestBuilder) WithTemp(temperature int) RequestBuilder {
	w.request.Params["temp"] = temperature
	return w
}

func (w requestBuilder) WithSpeed(speed int) RequestBuilder {
	w.request.Params["speed"] = speed
	return w
}

func (w requestBuilder) WithScene(scene Scene) RequestBuilder {
	w.request.Params["sceneId"] = scene
	return w
}

func (w requestBuilder) WithState(state bool) RequestBuilder {
	w.request.Params["state"] = state
	return w
}

func (w requestBuilder) Build() *Request {
	return w.request
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

type Response struct {
	Method string         `json:"method"`
	Env    string         `json:"env"`
	Result ResponseResult `json:"result"`
}

type ResponseResult struct {
	Mac     string `json:"mac"`
	Rssi    int    `json:"rssi"`
	State   bool   `json:"state"`
	SceneId int    `json:"sceneId"`
	Temp    int    `json:"temp"`
	Dimming int    `json:"dimming"`
}
