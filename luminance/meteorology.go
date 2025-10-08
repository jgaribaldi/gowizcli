package luminance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MeteorologyData struct {
	CloudCover    float64
	Precipitation float64
	Visibility    float64
	Thunderstorm  bool
	Elevation     float64
}

type Meteorology struct {
	baseUrl     string
	timeoutSecs int
}

func NewMeteorology(baseUrl string, timeoutSecs int) *Meteorology {
	return &Meteorology{
		baseUrl:     baseUrl,
		timeoutSecs: timeoutSecs,
	}
}

func (m Meteorology) GetCurrent(latitude, longitude float64) (*MeteorologyData, error) {
	q := url.Values{}
	q.Set("latitude", fmt.Sprintf("%v", latitude))
	q.Set("longitude", fmt.Sprintf("%v", longitude))
	q.Set("current", "cloud_cover,precipitation,visibility,weather_code")

	url := fmt.Sprintf("%s?%s", m.baseUrl, q.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	client := http.Client{Timeout: time.Duration(m.timeoutSecs) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting meteorology data: %d", res.StatusCode)
	}

	var out omApiResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &MeteorologyData{
		CloudCover:    out.Current.CloudCover,
		Precipitation: out.Current.Precipitation,
		Visibility:    out.Current.Visibility,
		Thunderstorm:  isThunderstorm(out.Current.WeatherCode),
		Elevation:     out.Elevation,
	}, nil
}

func isThunderstorm(code int) bool {
	return code == 95 || code == 96 || code == 99
}

type omCurrent struct {
	Time          string  `json:"time"`
	CloudCover    float64 `json:"cloud_cover"`
	Precipitation float64 `json:"precipitation"`
	Visibility    float64 `json:"visibility"`
	WeatherCode   int     `json:"weather_code"`
}

type omApiResponse struct {
	Elevation float64   `json:"elevation"`
	Current   omCurrent `json:"current"`
}
