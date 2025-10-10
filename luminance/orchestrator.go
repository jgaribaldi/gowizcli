package luminance

import "time"

type Orchestrator struct {
	solarElevationGetter    func(float64, float64) (*AstronomyData, error)
	weatherConditionsGetter func(float64, float64) (*MeteorologyData, error)
}

func NewOrchestrator(
	solarElevationGetter func(float64, float64) (*AstronomyData, error),
	weatherConditionsGetter func(float64, float64) (*MeteorologyData, error),
) *Orchestrator {
	return &Orchestrator{
		solarElevationGetter:    solarElevationGetter,
		weatherConditionsGetter: weatherConditionsGetter,
	}
}

func (o Orchestrator) GetCurrentLuminance(latitude, longitude float64) (float64, error) {
	astronomyData, err := o.solarElevationGetter(latitude, longitude)
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
