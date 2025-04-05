package utils

import "github.com/toxyl/math"

// Constants for RGB8 color conversion
const (
	RGB8Max = 255.0
)

// RGBToRGB8 converts RGB values to RGB8 values.
// RGB values should be in range [0,1].
// Returns RGB8 values in range [0,255].
func RGBToRGB8(r, g, b float64) (r8, g8, b8 float64) {
	// Clamp and convert to 8-bit values
	r8 = math.Round(r * RGB8Max)
	g8 = math.Round(g * RGB8Max)
	b8 = math.Round(b * RGB8Max)

	return r8, g8, b8
}

// RGB8ToRGB converts RGB8 values to RGB values.
// RGB8 values should be in range [0,255].
// Returns RGB values in range [0,1].
func RGB8ToRGB(r8, g8, b8 float64) (r, g, b float64) {
	// Clamp and convert to normalized values
	r = math.Clamp(r8, 0, RGB8Max) / RGB8Max
	g = math.Clamp(g8, 0, RGB8Max) / RGB8Max
	b = math.Clamp(b8, 0, RGB8Max) / RGB8Max

	return r, g, b
}
