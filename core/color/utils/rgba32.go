package utils

import "github.com/toxyl/math"

// Constants for RGBA32 color conversion
const (
	RGBA32Max = 4294967295.0
)

// RGBToRGBA32 converts RGB values to RGBA32 values.
// RGB values should be in range [0,1].
// Returns RGBA32 values in range [0,4294967295].
func RGBToRGBA32(r, g, b float64) (r32, g32, b32 float64) {
	// Clamp and convert to 32-bit values
	r32 = math.Round(r * RGBA32Max)
	g32 = math.Round(g * RGBA32Max)
	b32 = math.Round(b * RGBA32Max)

	return r32, g32, b32
}

// RGBA32ToRGB converts RGBA32 values to RGB values.
// RGBA32 values should be in range [0,4294967295].
// Returns RGB values in range [0,1].
func RGBA32ToRGB(r32, g32, b32 float64) (r, g, b float64) {
	// Clamp and convert to normalized values
	r = math.Clamp(r32, 0, RGBA32Max) / RGBA32Max
	g = math.Clamp(g32, 0, RGBA32Max) / RGBA32Max
	b = math.Clamp(b32, 0, RGBA32Max) / RGBA32Max

	return r, g, b
}
