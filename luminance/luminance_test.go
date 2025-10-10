package luminance

import (
	"math"
	"testing"
)

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

func TestEstimateLux_DayLight(t *testing.T) {
	type testCase struct {
		name     string
		in       ModelInput
		expected ModelOutput
	}

	tests := []testCase{
		{
			name:     "Clear sky at 30 degrees solar elevation, Linke turbidity 3, day of year 150, altitude 0 meters",
			in:       dayInput(30, 0, 0, 3, 150),
			expected: ModelOutput{Lux: 62947.9},
		},
		{
			name:     "Clear sky at 45 degrees solar elevation, Linke turbidity 3, day of year 100, altitude 0 meters",
			in:       dayInput(45, 0, 0, 3, 100),
			expected: ModelOutput{Lux: 92773.1},
		},
		{
			name:     "Overcast 100 percent at 45 degrees solar elevation, Linke turbidity 3, day of year 100, altitude 0 meters",
			in:       dayInput(45, 100, 0, 3, 100),
			expected: ModelOutput{Lux: 23193.3},
		},
		{
			name:     "Clear sky at 45 degrees solar elevation, Linke turbidity 2, day of year 200, altitude 0 meters",
			in:       dayInput(45, 0, 0, 3, 200),
			expected: ModelOutput{Lux: 90296.6},
		},
		{
			name:     "Clear sky at 45 degrees solar elevation, Linke turbidity 6, day of year 200, altitude 0 meters",
			in:       dayInput(45, 0, 0, 6, 200),
			expected: ModelOutput{Lux: 85227.0},
		},
		{
			name:     "Clear sky at 45 degrees solar elevation, Linke turbidity 3, day of year 120, altitude 2000 meters",
			in:       dayInput(45, 0, 2000, 3, 120),
			expected: ModelOutput{Lux: 94379.0},
		},
		{
			name:     "Clear sky at 45 degrees solar elevation, Linke turbidity 3, day of year 120, altitude 0 meters",
			in:       dayInput(45, 0, 0, 3, 120),
			expected: ModelOutput{Lux: 91773.5},
		},
		{
			name:     "Clear sky at 1 degree solar elevation near the horizon, Linke turbidity 3, day of year 80, altitude 0 meters",
			in:       dayInput(1.0, 0, 0, 3, 80),
			expected: ModelOutput{Lux: 1077.6},
		},
		{
			name:     "Clear sky at 60 degrees solar elevation (high Sun), Linke turbidity 3, day of year 150, altitude 0 meters",
			in:       dayInput(60, 0, 0, 3, 150),
			expected: ModelOutput{Lux: 111880.6},
		},
		{
			name:     "Clear sky at 30 degrees solar elevation near perihelion, Linke turbidity 3, day of year 3, altitude 0 meters",
			in:       dayInput(30, 0, 0, 3, 3),
			expected: ModelOutput{Lux: 66893.3},
		},
		{
			name:     "Clear sky at 30 degrees solar elevation near aphelion, Linke turbidity 3, day of year 185, altitude 0 meters",
			in:       dayInput(30, 0, 0, 3, 185),
			expected: ModelOutput{Lux: 62624.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, _ := EstimateLux(tt.in)
			gotRounded := math.Round(out.Lux*10) / 10
			if gotRounded != tt.expected.Lux {
				t.Fatalf("got %f; expected %f", gotRounded, tt.expected.Lux)
			}
		})
	}
}

func dayInput(solarDeg, cloudPct, altitudeM, linkeTL float64, dayOfYear int) ModelInput {
	return ModelInput{
		SolarElevationDeg:    solarDeg,
		CloudCoverPercentage: cloudPct,
		AltitudeMeters:       altitudeM,
		DayOfYear:            dayOfYear,
		LinkeTurbidity:       linkeTL,
	}
}
