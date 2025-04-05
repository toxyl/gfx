package utils

import "github.com/toxyl/math"

// RGBToCMY converts RGB to CMY color space.
// RGB values should be in range [0,1].
// Returns CMY values in range [0,1].
func RGBToCMY(r, g, b float64) (c, m, y float64) {
	// Convert RGB to CMY
	c = 1 - r
	m = 1 - g
	y = 1 - b

	return c, m, y
}

// CMYToRGB converts CMY to RGB color space.
// CMY values should be in range [0,1].
// Returns RGB values in range [0,1].
func CMYToRGB(c, m, y float64) (r, g, b float64) {
	// Convert CMY to RGB
	r = 1 - c
	g = 1 - m
	b = 1 - y

	// Clamp values to [0,1]
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return r, g, b
}
