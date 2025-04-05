package projections

import (
	"math"
)

// RectPolar implements the rectangular to polar projection
type RectPolar struct {
	BaseProjection
}

// NewRectPolar creates a new rectangular to polar projection
func NewRectPolar() *RectPolar {
	return &RectPolar{}
}

// To converts from cartesian coordinates to polar coordinates
func (p *RectPolar) To(x, y, width, height float64) (radius, angle float64) {
	// Handle edge cases
	if x == width/2 && y == height/2 {
		return 0, 0
	}

	// Convert to normalized coordinates centered at origin
	xn := 2*x/width - 1
	yn := 1 - 2*y/height

	// Calculate radius and angle
	radius = math.Sqrt(xn*xn + yn*yn)
	angle = math.Atan2(yn, xn)

	// Normalize angle to [0, 2π]
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Handle edge cases for radius
	if x == width-1 && y == height/2 {
		radius = 1
		angle = 0
	} else if x == 0 && y == height/2 {
		radius = 1
		angle = math.Pi
	} else if x == width/2 && y == 0 {
		radius = 1
		angle = math.Pi / 2
	} else if x == width/2 && y == height-1 {
		radius = 1
		angle = 3 * math.Pi / 2
	}

	return radius, angle
}

// From converts from polar coordinates to cartesian coordinates
func (p *RectPolar) From(radius, angle, width, height float64) (x, y float64) {
	// Handle edge cases
	if radius == 0 {
		return width / 2, height / 2
	}

	// Normalize angle to [0, 2π]
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Handle edge cases for radius = 1
	if radius == 1 {
		switch angle {
		case 0:
			return width - 1, height / 2
		case math.Pi:
			return 0, height / 2
		case math.Pi / 2:
			return width / 2, 0
		case 3 * math.Pi / 2:
			return width / 2, height - 1
		}
	}

	// Convert to cartesian coordinates
	xn := radius * math.Cos(angle)
	yn := radius * math.Sin(angle)

	// Convert to pixel coordinates
	x = (xn + 1) / 2 * width
	y = (1 - yn) / 2 * height

	return x, y
}
