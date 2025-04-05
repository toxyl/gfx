package utils

import "github.com/toxyl/math"

// RGBToYUV converts RGB to YUV color space.
// Uses BT.601 standard for conversion.
func RGBToYUV(r, g, b float64) (y, u, v float64) {
	// BT.601 conversion
	y = 0.299*r + 0.587*g + 0.114*b
	u = -0.14713*r - 0.28886*g + 0.436*b
	v = 0.615*r - 0.51499*g - 0.10001*b

	return y, u, v
}

// YUVToRGB converts YUV to RGB color space.
// Uses BT.601 standard for conversion.
func YUVToRGB(y, u, v float64) (r, g, b float64) {
	// BT.601 conversion
	r = y + 1.13983*v
	g = y - 0.39465*u - 0.58060*v
	b = y + 2.03211*u

	// Clamp values to [0,1]
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return r, g, b
}
