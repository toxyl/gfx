package projections

import (
	"math"
)

// Stereographic implements the stereographic projection
type Stereographic struct {
	BaseProjection
}

// NewStereographic creates a new stereographic projection
func NewStereographic() *Stereographic {
	return &Stereographic{}
}

// To converts from geographic coordinates to cartesian coordinates
func (p *Stereographic) To(latitude, longitude, width, height float64) (x, y float64) {
	// Handle edge cases
	if latitude == 90 {
		return width / 2, 0
	} else if latitude == -90 {
		return width / 2, height - 1
	}

	// Convert to radians
	latRad := latitude * math.Pi / 180.0
	lonRad := longitude * math.Pi / 180.0

	// Stereographic projection
	k := 2.0 / (1 + math.Sin(latRad)*math.Sin(0) + math.Cos(latRad)*math.Cos(0)*math.Cos(lonRad))
	x = width/2 + k*math.Cos(latRad)*math.Sin(lonRad)*width/4
	y = height/2 - k*(math.Cos(0)*math.Sin(latRad)-math.Sin(0)*math.Cos(latRad)*math.Cos(lonRad))*height/4

	// Handle longitude edge cases
	if longitude == 180 {
		x = width - 1
	} else if longitude == -180 {
		x = 0
	}

	return x, y
}

// From converts from cartesian coordinates to geographic coordinates
func (p *Stereographic) From(x, y, width, height float64) (latitude, longitude float64) {
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
	xn := 4 * (x - width/2) / width
	yn := 4 * (height/2 - y) / height
	rho := math.Sqrt(xn*xn + yn*yn)
	c := 2 * math.Atan2(rho, 2)

	// Convert to geographic coordinates
	latRad := math.Asin(math.Cos(c)*math.Sin(0) + yn*math.Sin(c)*math.Cos(0)/rho)
	lonRad := math.Atan2(xn*math.Sin(c), rho*math.Cos(0)*math.Cos(c)-yn*math.Sin(0)*math.Sin(c))

	latitude = latRad * 180.0 / math.Pi
	longitude = lonRad * 180.0 / math.Pi

	return latitude, longitude
}
