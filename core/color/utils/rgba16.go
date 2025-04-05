package utils

import "github.com/toxyl/math"

// Constants for RGBA16 color conversion
const (
	RGBA16Max = 65535.0
)

// RGBToRGBA16 converts RGB values to RGBA16 values.
// RGB values should be in range [0,1].
// Returns RGBA16 values in range [0,65535].
func RGBToRGBA16(r, g, b float64) (r16, g16, b16 float64) {
	// Clamp and convert to 16-bit values
	r16 = math.Round(r * RGBA16Max)
	g16 = math.Round(g * RGBA16Max)
	b16 = math.Round(b * RGBA16Max)

	return r16, g16, b16
}

// RGBA16ToRGB converts RGBA16 values to RGB values.
// RGBA16 values should be in range [0,65535].
// Returns RGB values in range [0,1].
func RGBA16ToRGB(r16, g16, b16 float64) (r, g, b float64) {
	// Clamp and convert to normalized values
	r = math.Clamp(r16, 0, RGBA16Max) / RGBA16Max
	g = math.Clamp(g16, 0, RGBA16Max) / RGBA16Max
	b = math.Clamp(b16, 0, RGBA16Max) / RGBA16Max

	return r, g, b
}
