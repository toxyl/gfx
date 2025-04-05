package utils

import "github.com/toxyl/math"

// Constants for RGBA64 color conversion
const (
	RGBA64Max = 18446744073709551615.0
)

// RGBToRGBA64 converts RGB values to RGBA64 values.
// RGB values should be in range [0,1].
// Returns RGBA64 values in range [0,18446744073709551615].
func RGBToRGBA64(r, g, b float64) (r64, g64, b64 float64) {
	// Clamp and convert to 64-bit values
	r64 = math.Round(r * RGBA64Max)
	g64 = math.Round(g * RGBA64Max)
	b64 = math.Round(b * RGBA64Max)

	return r64, g64, b64
}

// RGBA64ToRGB converts RGBA64 values to RGB values.
// RGBA64 values should be in range [0,18446744073709551615].
// Returns RGB values in range [0,1].
func RGBA64ToRGB(r64, g64, b64 float64) (r, g, b float64) {
	// Clamp and convert to normalized values
	r = math.Clamp(r64, 0, RGBA64Max) / RGBA64Max
	g = math.Clamp(g64, 0, RGBA64Max) / RGBA64Max
	b = math.Clamp(b64, 0, RGBA64Max) / RGBA64Max

	return r, g, b
}
