package utils

import "github.com/toxyl/math"

// RGBToCMYK converts RGB to CMYK color space.
// All values are in range [0,1].
func RGBToCMYK(r, g, b float64) (c, m, y, k float64) {
	// Calculate key (black) component
	k = 1 - math.Max(r, math.Max(g, b))

	// Calculate CMY components
	if k < 1 {
		c = (1 - r - k) / (1 - k)
		m = (1 - g - k) / (1 - k)
		y = (1 - b - k) / (1 - k)
	} else {
		c = 0
		m = 0
		y = 0
	}

	return c, m, y, k
}

// CMYKToRGB converts CMYK to RGB color space.
// All values are in range [0,1].
func CMYKToRGB(c, m, y, k float64) (r, g, b float64) {
	// Convert CMYK to RGB
	r = (1 - c) * (1 - k)
	g = (1 - m) * (1 - k)
	b = (1 - y) * (1 - k)

	// Clamp values to [0,1]
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return r, g, b
}
