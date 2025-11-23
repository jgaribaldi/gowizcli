package luminance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Astronomy interface {
	GetSolarElevation(latitude, longitude float64) (*AstronomyData, error)
}

type AstronomyData struct {
	SunAltitude float64 `json:"sun_altitude"`
}

type IpGeolocation struct {
	baseUrl     string
	apiKey      string
	timeoutSecs int
}

func NewIpGeolocation(baseUrl string, apiKey string, timeoutSecs int) *IpGeolocation {
	return &IpGeolocation{
		baseUrl:     baseUrl,
		apiKey:      apiKey,
		timeoutSecs: timeoutSecs,
	}
}

func (i IpGeolocation) GetSolarElevation(latitude, longitude float64) (*AstronomyData, error) {
	strLat := strconv.FormatFloat(latitude, 'f', -1, 64)
	strLong := strconv.FormatFloat(longitude, 'f', -1, 64)

	q := url.Values{}
	q.Set("apiKey", i.apiKey)
	q.Set("lat", strLat)
	q.Set("long", strLong)

	url := fmt.Sprintf("%s?%s", i.baseUrl, q.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	client := http.Client{Timeout: time.Duration(i.timeoutSecs) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting data from ipgeolocation: %d", res.StatusCode)
	}

	var out ipgApiResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out.Astronomy, nil
}

type ipgApiResponse struct {
	Astronomy AstronomyData `json:"astronomy"`
}
