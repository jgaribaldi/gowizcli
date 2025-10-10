package luminance

import "testing"

func TestEstimateLux_LuxZero(t *testing.T) {
	input := ModelInput{
		SolarElevationDeg:    10.0,
		CloudCoverPercentage: 50.0,
		AltitudeMeters:       100.0,
		DayOfYear:            150,
		LinkeTurbidity:       3.0,
	}

	out, _ := EstimateLux(input)

	if out.Lux <= 0 {
		t.Errorf("Expected Lux == 0, got: %v", out.Lux)
	}
}
