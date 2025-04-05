package projections

import (
	"math"
	"testing"
)

const (
	testWidth  = 800.0
	testHeight = 600.0
	epsilon    = 1e-10
)

func TestEquirectangular(t *testing.T) {
	proj := NewEquirectangular()

	// Test known points
	tests := []struct {
		name      string
		lat, lon  float64
		expectedX float64
		expectedY float64
	}{
		{"North Pole", 90, 0, testWidth / 2, 0},
		{"South Pole", -90, 0, testWidth / 2, testHeight - 1},
		{"Equator Prime Meridian", 0, 0, testWidth / 2, testHeight / 2},
		{"Equator 180°", 0, 180, testWidth - 1, testHeight / 2},
		{"Equator -180°", 0, -180, 0, testHeight / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := proj.To(tt.lat, tt.lon, testWidth, testHeight)
			if math.Abs(x-tt.expectedX) > epsilon || math.Abs(y-tt.expectedY) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}

			// Test round trip
			lat, lon := proj.From(x, y, testWidth, testHeight)
			if math.Abs(lat-tt.lat) > epsilon || math.Abs(lon-tt.lon) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", lat, lon, tt.lat, tt.lon)
			}
		})
	}
}

func TestMercator(t *testing.T) {
	proj := NewMercator()

	// Test known points
	tests := []struct {
		name      string
		lat, lon  float64
		expectedX float64
		expectedY float64
	}{
		{"Equator Prime Meridian", 0, 0, testWidth / 2, testHeight / 2},
		{"Equator 180°", 0, 180, testWidth - 1, testHeight / 2},
		{"Equator -180°", 0, -180, 0, testHeight / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := proj.To(tt.lat, tt.lon, testWidth, testHeight)
			if math.Abs(x-tt.expectedX) > epsilon || math.Abs(y-tt.expectedY) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}

			// Test round trip
			lat, lon := proj.From(x, y, testWidth, testHeight)
			if math.Abs(lat-tt.lat) > epsilon || math.Abs(lon-tt.lon) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", lat, lon, tt.lat, tt.lon)
			}
		})
	}
}

func TestSinusoidal(t *testing.T) {
	proj := NewSinusoidal()

	// Test known points
	tests := []struct {
		name      string
		lat, lon  float64
		expectedX float64
		expectedY float64
	}{
		{"Equator Prime Meridian", 0, 0, testWidth / 2, testHeight / 2},
		{"Equator 180°", 0, 180, testWidth - 1, testHeight / 2},
		{"Equator -180°", 0, -180, 0, testHeight / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := proj.To(tt.lat, tt.lon, testWidth, testHeight)
			if math.Abs(x-tt.expectedX) > epsilon || math.Abs(y-tt.expectedY) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}

			// Test round trip
			lat, lon := proj.From(x, y, testWidth, testHeight)
			if math.Abs(lat-tt.lat) > epsilon || math.Abs(lon-tt.lon) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", lat, lon, tt.lat, tt.lon)
			}
		})
	}
}

func TestStereographic(t *testing.T) {
	proj := NewStereographic()

	// Test known points
	tests := []struct {
		name      string
		lat, lon  float64
		expectedX float64
		expectedY float64
	}{
		{"Equator Prime Meridian", 0, 0, testWidth / 2, testHeight / 2},
		{"Equator 180°", 0, 180, testWidth - 1, testHeight / 2},
		{"Equator -180°", 0, -180, 0, testHeight / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := proj.To(tt.lat, tt.lon, testWidth, testHeight)
			if math.Abs(x-tt.expectedX) > epsilon || math.Abs(y-tt.expectedY) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}

			// Test round trip
			lat, lon := proj.From(x, y, testWidth, testHeight)
			if math.Abs(lat-tt.lat) > epsilon || math.Abs(lon-tt.lon) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", lat, lon, tt.lat, tt.lon)
			}
		})
	}
}

func TestRectPolar(t *testing.T) {
	proj := NewRectPolar()

	// Test known points
	tests := []struct {
		name          string
		x, y          float64
		expectedR     float64
		expectedTheta float64
	}{
		{"Center", testWidth / 2, testHeight / 2, 0, 0},
		{"Right", testWidth - 1, testHeight / 2, 1, 0},
		{"Left", 0, testHeight / 2, 1, math.Pi},
		{"Top", testWidth / 2, 0, 1, math.Pi / 2},
		{"Bottom", testWidth / 2, testHeight - 1, 1, 3 * math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, theta := proj.To(tt.x, tt.y, testWidth, testHeight)
			if math.Abs(r-tt.expectedR) > epsilon || math.Abs(theta-tt.expectedTheta) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", r, theta, tt.expectedR, tt.expectedTheta)
			}

			// Test round trip
			x, y := proj.From(r, theta, testWidth, testHeight)
			if math.Abs(x-tt.x) > epsilon || math.Abs(y-tt.y) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", x, y, tt.x, tt.y)
			}
		})
	}
}

func TestPolarRect(t *testing.T) {
	proj := NewPolarRect()

	// Test known points
	tests := []struct {
		name      string
		r, theta  float64
		expectedX float64
		expectedY float64
	}{
		{"Center", 0, 0, testWidth / 2, testHeight / 2},
		{"Right", 1, 0, testWidth - 1, testHeight / 2},
		{"Left", 1, math.Pi, 0, testHeight / 2},
		{"Top", 1, math.Pi / 2, testWidth / 2, 0},
		{"Bottom", 1, 3 * math.Pi / 2, testWidth / 2, testHeight - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := proj.To(tt.r, tt.theta, testWidth, testHeight)
			if math.Abs(x-tt.expectedX) > epsilon || math.Abs(y-tt.expectedY) > epsilon {
				t.Errorf("To() = (%v, %v), want (%v, %v)", x, y, tt.expectedX, tt.expectedY)
			}

			// Test round trip
			r, theta := proj.From(x, y, testWidth, testHeight)
			if math.Abs(r-tt.r) > epsilon || math.Abs(theta-tt.theta) > epsilon {
				t.Errorf("From() = (%v, %v), want (%v, %v)", r, theta, tt.r, tt.theta)
			}
		})
	}
}
