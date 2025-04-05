package utils

import "github.com/toxyl/math"

// RGBToYCbCr converts RGB to YCbCr color space.
// Uses BT.601 standard for conversion.
func RGBToYCbCr(r, g, b float64) (y, cb, cr float64) {
	// BT.601 conversion
	y = 0.299*r + 0.587*g + 0.114*b
	cb = -0.168736*r - 0.331264*g + 0.5*b
	cr = 0.5*r - 0.418688*g - 0.081312*b

	return y, cb, cr
}

// YCbCrToRGB converts YCbCr to RGB color space.
// Uses BT.601 standard for conversion.
func YCbCrToRGB(y, cb, cr float64) (r, g, b float64) {
	// BT.601 conversion
	r = y + 1.402*cr
	g = y - 0.344136*cb - 0.714136*cr
	b = y + 1.772*cb

	// Clamp values to [0,1]
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return r, g, b
}
