package projections

import (
	"math"
)

// Mercator implements the mercator projection
type Mercator struct {
	BaseProjection
}

// NewMercator creates a new mercator projection
func NewMercator() *Mercator {
	return &Mercator{}
}

// To converts from geographic coordinates to cartesian coordinates
func (p *Mercator) To(latitude, longitude, width, height float64) (x, y float64) {
	// Handle edge cases
	if latitude == 90 {
		return width / 2, 0
	} else if latitude == -90 {
		return width / 2, height - 1
	}

	// Convert to radians
	latRad := latitude * math.Pi / 180.0
	lonRad := longitude * math.Pi / 180.0

	// Mercator projection
	x = (lonRad + math.Pi) / (2 * math.Pi) * width
	y = (1 - math.Log(math.Tan(latRad/2+math.Pi/4))/math.Pi) / 2 * height

	// Handle longitude edge cases
	if longitude == 180 {
		x = width - 1
	} else if longitude == -180 {
		x = 0
	}

	return x, y
}

// From converts from cartesian coordinates to geographic coordinates
func (p *Mercator) From(x, y, width, height float64) (latitude, longitude float64) {
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
	yn := 2*y/height - 1
	lonRad := (2*x/width - 1) * math.Pi

	// Convert to geographic coordinates
	latRad := 2*math.Atan(math.Exp(-yn*math.Pi)) - math.Pi/2
	longitude = lonRad * 180.0 / math.Pi
	latitude = latRad * 180.0 / math.Pi

	return latitude, longitude
}
