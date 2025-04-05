package projections

import (
	"math"
)

// Sinusoidal implements the sinusoidal projection
type Sinusoidal struct {
	BaseProjection
}

// NewSinusoidal creates a new sinusoidal projection
func NewSinusoidal() *Sinusoidal {
	return &Sinusoidal{}
}

// To converts from geographic coordinates to cartesian coordinates
func (p *Sinusoidal) To(latitude, longitude, width, height float64) (x, y float64) {
	// Handle edge cases
	if latitude == 90 {
		return width / 2, 0
	} else if latitude == -90 {
		return width / 2, height - 1
	}

	// Convert to radians
	latRad := latitude * math.Pi / 180.0
	lonRad := longitude * math.Pi / 180.0

	// Sinusoidal projection
	x = (lonRad*math.Cos(latRad) + math.Pi) / (2 * math.Pi) * width
	y = (math.Pi/2 - latRad) / math.Pi * height

	// Handle longitude edge cases
	if longitude == 180 {
		x = width - 1
	} else if longitude == -180 {
		x = 0
	}

	return x, y
}

// From converts from cartesian coordinates to geographic coordinates
func (p *Sinusoidal) From(x, y, width, height float64) (latitude, longitude float64) {
	// Handle edge cases
	if y == 0 {
		return 90, 0
	} else if y == height-1 {
		return -90, 0
	}

	if x == width-1 {
		return 0, 180
	} else if x == 0 {
		return 0, -180
	}

	// Convert to normalized coordinates
	latRad := math.Pi/2 - (y/height)*math.Pi
	lonRad := ((x/width)*2*math.Pi - math.Pi) / math.Cos(latRad)

	// Convert to degrees
	latitude = latRad * 180.0 / math.Pi
	longitude = lonRad * 180.0 / math.Pi

	return latitude, longitude
}
