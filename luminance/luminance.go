package luminance

import "time"

type Luminance struct {
	astronomy   Astronomy
	meteorology Meteorology
}

func NewLuminance(astronomy Astronomy, meteorology Meteorology) Luminance {
	return Luminance{
		astronomy:   astronomy,
		meteorology: meteorology,
	}
}

func (l Luminance) GetCurrent(latitude, longitude float64) (float64, error) {
	astronomyData, err := l.astronomy.GetSolarElevation(latitude, longitude)
	if err != nil {
		return -1.0, err
	}

	meteorologyData, err := l.meteorology.GetCurrent(latitude, longitude)
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
