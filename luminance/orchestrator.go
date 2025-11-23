package luminance

import "time"

type Luminance struct {
	astronomy               Astronomy
	weatherConditionsGetter func(float64, float64) (*MeteorologyData, error)
}

func NewLuminance(
	astronomy Astronomy,
	weatherConditionsGetter func(float64, float64) (*MeteorologyData, error),
) *Luminance {
	return &Luminance{
		astronomy:               astronomy,
		weatherConditionsGetter: weatherConditionsGetter,
	}
}

func (o Luminance) GetCurrent(latitude, longitude float64) (float64, error) {
	astronomyData, err := o.astronomy.GetSolarElevation(latitude, longitude)
	if err != nil {
		return -1.0, err
	}

	meteorologyData, err := o.weatherConditionsGetter(latitude, longitude)
	if err != nil {
		return -1.0, err
	}

	modelInput := ModelInput{
		SolarElevationDeg:    astronomyData.SunAltitude,
		CloudCoverPercentage: meteorologyData.CloudCover,
		AltitudeMeters:       meteorologyData.Elevation,
		DayOfYear:            time.Time.YearDay(time.Now()),
		LinkeTurbidity:       LinkeTurbidityDefault,
	}
	luminance := EstimateLux(modelInput)

	return luminance.Lux, nil
}
