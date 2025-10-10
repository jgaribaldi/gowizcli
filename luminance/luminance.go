package luminance

import "math"

type ModelInput struct {
	SolarElevationDeg    float64
	CloudCoverPercentage float64
	AltitudeMeters       float64
	DayOfYear            int
	LinkeTurbidity       float64
}

type ModelOutput struct {
	Lux float64
}

func EstimateLux(input ModelInput) (*ModelOutput, error) {
	if isDaytime(input.SolarElevationDeg) {
		result := luxForDaytime(input)
		return &result, nil
	}

	return &ModelOutput{
		0,
	}, nil
}

func isDaytime(solarElevationDeg float64) bool {
	return solarElevationDeg > 0
}

func luxForDaytime(input ModelInput) ModelOutput {
	airMass := airmassKastenYoung(input.SolarElevationDeg)
	e0 := eccentricityFactor(input.DayOfYear)
	ghiClear := clearSkyGHI(input.AltitudeMeters, LinkeTurbidityDefault, input.SolarElevationDeg, airMass, e0)
	kc := cloudClearSkyIndex(input.CloudCoverPercentage)

	lux := ghiClear * kc * LuminousEfficacy_LuxPerWm2
	return ModelOutput{
		Lux: lux,
	}
}

func airmassKastenYoung(solarElevationDeg float64) float64 {
	solarElevationRad := radians(solarElevationDeg)
	return 1.0 / (math.Sin(solarElevationRad) * AirMassKastenYoungA * math.Pow(AirMassKastenYoungB+(90.0-solarElevationDeg), -AirMassKastenYoungC))
}

func eccentricityFactor(dayOfYear int) float64 {
	return 1.0 + 0.033*math.Cos(2.0*math.Pi*float64(dayOfYear)/365.0)
}

func clearSkyGHI(altitudeM, linkeTurbidity, solarElevationDeg, airMass, eccentricityFactor float64) float64 {
	solarElevationRad := radians(solarElevationDeg)
	io := SolarIrradianceExtraterrestrial_Wm2 * eccentricityFactor
	fh1 := AltitudeTransmittanceBase + AltitudeTransmittanceSlopePerMeter*altitudeM
	fh2 := TurbidityScalingBase

	return ClearSkyModelPrefactor_IneichenKasten * io * math.Sin(solarElevationRad) * math.Exp(-ClearSkyOpticalDepthCoeff*airMass*(fh1+fh2*(linkeTurbidity-1.0)))
}

func cloudClearSkyIndex(cloudCoverPercentage float64) float64 {
	cloudFraction := cloudCoverPercentage / 100.0
	if cloudFraction <= 0.0 {
		return 1.0
	}
	if cloudFraction >= 1.0 {
		return 1.0 - CloudClearSkyLossCoefficient
	}

	return 1.0 - CloudClearSkyLossCoefficient*math.Pow(cloudFraction, CloudClearSkyExponent)
}

const SolarIrradianceExtraterrestrial_Wm2 = 1366.1

const ClearSkyOpticalDepthCoeff = 0.027

const (
	AltitudeTransmittanceBase          = 0.128
	AltitudeTransmittanceSlopePerMeter = -0.000367
	TurbidityScalingBase               = 0.505
)

const (
	AirMassKastenYoungA = 0.50572
	AirMassKastenYoungB = 6.07995
	AirMassKastenYoungC = 1.6364
)

const (
	CloudClearSkyLossCoefficient = 0.75
	CloudClearSkyExponent        = 3.4
)

const ClearSkyModelPrefactor_IneichenKasten = 0.84

const LuminousEfficacy_LuxPerWm2 = 120.0
