package luminance

import (
	"math"
)

type Luminance struct {
	ipGeolocation *IpGeolocation
	meteorology   *Meteorology
}

func NewLuminance(ipGeolocation *IpGeolocation, meteorology *Meteorology) *Luminance {
	return &Luminance{
		ipGeolocation: ipGeolocation,
		meteorology:   meteorology,
	}
}

func (l Luminance) CalculateOutsideLuminance() (float64, error) {
	data, err := l.meteorology.GetCurrent(-34.60734, -58.44329)
	if err != nil {
		return 0, err
	}

	astronomyData, err := l.ipGeolocation.GetSolarElevation("Buenos Aires, ARG")
	if err != nil {
		return 0, err
	}

	contextFactors := contextFactors{
		isUrban:      true,
		visibilityKm: data.Visibility / 1000.0,
		altitudeM:    data.Elevation,
	}
	luminance := estimateLuminanceLux(
		astronomyData.SunAltitude,
		data.CloudCover,
		data.Precipitation,
		data.Thunderstorm,
		contextFactors,
	)
	return luminance, nil
}

type contextFactors struct {
	isUrban      bool
	altitudeM    float64
	visibilityKm float64
}

func estimateLuminanceLux(
	sunElevationDeg float64,
	cloudCoverPct float64,
	precipitMmPerHour float64,
	thunderstorm bool,
	contextFactors contextFactors,
) float64 {
	if sunElevationDeg <= 0.0 {
		return 0.0
	}

	sunGeometryTerm := math.Sin(radians(sunElevationDeg))
	valueLux := clearSkyPeakLux * sunGeometryTerm * cloudTransmittance(cloudCoverPct)
	valueLux = valueLux * rainFactor(precipitMmPerHour)

	if thunderstorm {
		valueLux = valueLux * thunderstormFactor
	}

	valueLux = valueLux * visibilityFactorKm(contextFactors.visibilityKm)
	valueLux = valueLux * urbanFactor(contextFactors.isUrban)
	valueLux = valueLux * altitudeBoost(contextFactors.altitudeM)

	return max(0.0, valueLux)
}

func cloudTransmittance(cloudCoverPct float64) float64 {
	cloudFraction := max(0.0, min(100, cloudCoverPct)) / 100.0
	return 1.0 - cloudAttenuationCoef*math.Pow(cloudFraction, cloudAttenuationExponent)
}

func rainFactor(precipitMmPerHour float64) float64 {
	if precipitMmPerHour <= 0.0 {
		return 1.0
	}
	effective := min(rainCapMmPerHour, precipitMmPerHour)
	candidate := 1.0 - rainSlopePerMm*effective
	return max(rainMinFactor, candidate)
}

func visibilityFactorKm(visibilityKm float64) float64 {
	for _, vb := range visibilityBandsKm {
		if visibilityKm < vb.Threshold {
			return vb.MultiplicationFactor
		}
	}
	return 1.0
}

func urbanFactor(isUrban bool) float64 {
	if isUrban {
		return 1.0 - UrbanReductionFraction
	}
	return 1.0
}

func altitudeBoost(altitudeM float64) float64 {
	if altitudeM <= 0.0 {
		return 1.0
	}

	boost := altitudeBoostPer100m * altitudeM / 100.0
	boost = min(boost, altitudeMaxBoost)
	return 1.0 + boost
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

const (
	clearSkyPeakLux          = 100000.0
	cloudAttenuationCoef     = 0.75
	cloudAttenuationExponent = 3.4
	rainMinFactor            = 0.3
	rainSlopePerMm           = 0.1
	rainCapMmPerHour         = 7.0
	thunderstormFactor       = 0.4
	UrbanReductionFraction   = 0.2
	altitudeBoostPer100m     = 0.01
	altitudeMaxBoost         = 0.2
)

type visibilityBandKm struct {
	Threshold            float64
	MultiplicationFactor float64
}

var visibilityBandsKm = []visibilityBandKm{{1.0, 0.3}, {5.0, 0.6}, {10.0, 0.8}}
