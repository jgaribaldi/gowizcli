package luminance

import "time"

type Luminance struct {
	Astronomy   Astronomy
	Meteorology Meteorology
}

func (l Luminance) GetCurrent(latitude, longitude float64) (float64, error) {
	astronomyData, err := l.Astronomy.GetSolarElevation(latitude, longitude)
	if err != nil {
		return -1.0, err
	}

	meteorologyData, err := l.Meteorology.GetCurrent(latitude, longitude)
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
