package luminance

import (
	"math"
	"time"
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

const (
	LinkeTurbidityDefault float64 = 3.0
)

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
		time.Time.YearDay(time.Now()),
		LinkeTurbidityDefault,
		data.CloudCover,
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
	dayOfYear int,
	linkeTurbidity float64,
	cloudCoverPct float64,
	contextFactors contextFactors,
) float64 {
	zenith := 90.0 - sunElevationDeg

	// clear sky
	ipcs := getIneichenPerezForClearSky(zenith, contextFactors.altitudeM, linkeTurbidity, dayOfYear)

	// cloud attenuation
	kc := getCloudTransmittanceKastenCzeplak(cloudCoverPct)

	dniAllSky := max(0.0, ipcs.dniCs*kc)
	dhiAllSky := max(0.0, ipcs.dhiCs*kc)

	evLux := irradianceToIlluminance(dniAllSky, dhiAllSky, sunElevationDeg)

	if sunElevationDeg <= 0.0 {
		evLux = max(evLux, twilightLux(sunElevationDeg))
	}

	return evLux
}

type IneichenPerezClearSky struct {
	ghiCs float64
	dniCs float64
	dhiCs float64
}

func getIneichenPerezForClearSky(
	solarZenithDeg float64,
	altitudeM float64,
	linkeTurbidity float64,
	dayOfYear int,
) IneichenPerezClearSky {
	cos := max(math.Cos(radians(solarZenithDeg)), 0.0)
	if cos <= 0.0 {
		return IneichenPerezClearSky{}
	}

	iExt := extraterrestrialIrradianceNormal(dayOfYear)
	am := airmassKastenYoung(solarZenithDeg)

	fh1 := math.Exp(-altitudeM / Fh1DenM)
	fh2 := math.Exp(-altitudeM / Fh2DenM)

	dniCs := BBeamAdjust * iExt * math.Exp(-BeamAttenCoef*am*(linkeTurbidity-1.0))
	ghiCs := A1 * cos * iExt * math.Exp(-A2*am*(fh1+fh2*(linkeTurbidity-1.0)))
	dhiCs := max(0.0, ghiCs-dniCs*cos)

	return IneichenPerezClearSky{
		dniCs: dhiCs,
		ghiCs: ghiCs,
		dhiCs: dhiCs,
	}
}

func extraterrestrialIrradianceNormal(dayOfYear int) float64 {
	return SolarConstantWM2 * excentricityCorrection(dayOfYear)
}

func excentricityCorrection(dayOfYear int) float64 {
	return 1.0 + ExcentricityCosCoef*math.Cos(2.0*math.Pi*float64(dayOfYear)/365.0)
}

func airmassKastenYoung(zenithDeg float64) float64 {
	z := max(0.0, min(zenithDeg, 90.0))
	denom := math.Cos(radians(z)) + KyA*math.Pow((KyBDeg-z), (-KyC))
	return 1.0 / denom
}

func getCloudTransmittanceKastenCzeplak(cloudCoverPercent float64) float64 {
	clamped := max(0.0, min(cloudCoverPercent, 100.0))
	oktas := math.Round(clamped / 12.5)
	kc := 1.0 - math.Pow(CCoef*(oktas/0.8), Exponent)
	return max(0.0, min(kc, 1.0))
}

const (
	SolarConstantWM2    float64 = 1367.0
	ExcentricityCosCoef float64 = 0.033
	KyA                 float64 = 0.50572
	KyBDeg              float64 = 96.07995
	KyC                 float64 = 1.6364
	A1                  float64 = (5.09e-5 + 0.868)
	A2                  float64 = (3.92e-5 + 0.0387)
	Fh1DenM             float64 = 8000.0
	Fh2DenM             float64 = 1250.0
	BBeamAdjust         float64 = 1.0
	BeamAttenCoef       float64 = 0.09
	CCoef               float64 = 0.75
	Exponent            float64 = 3.4
)

const (
	GlobalLuminancePerW  float64 = 120.0
	DiffuseLuminancePerW float64 = 130.0
	DirectLuminancePerW  float64 = 93.0
)

func irradianceToIlluminance(dniWm2 float64, dhiWm2 float64, solarZenithDeg float64) float64 {
	cosz := max(math.Cos(radians(solarZenithDeg)), 0.0)
	directComponent := DirectLuminancePerW * max(0.0, dniWm2) * cosz
	diffuseComponent := DiffuseLuminancePerW * max(0.0, dhiWm2)
	return directComponent + diffuseComponent
}

const (
	CivilTwilightStartDeg float64 = -6.0
	TwilightMaxLux        float64 = 200.0
	TwilightShapeK        float64 = 1.2
)

func twilightLux(sunElevationDeg float64) float64 {
	if sunElevationDeg <= CivilTwilightStartDeg-4.0 {
		return 0.0
	}

	x := sunElevationDeg - CivilTwilightStartDeg
	return TwilightMaxLux / (1.0 + math.Exp(-TwilightShapeK*x))
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
