package utils

import "github.com/toxyl/math"

// RGBToYIQ converts RGB to YIQ color space.
// Uses NTSC standard for conversion.
// RGB values should be in range [0,1].
// Returns Y in range [0,1], I and Q in range [-0.5,0.5].
func RGBToYIQ(r, g, b float64) (y, i, q float64) {
	// NTSC YIQ conversion matrix
	y = 0.299*r + 0.587*g + 0.114*b
	i = 0.596*r - 0.275*g - 0.321*b
	q = 0.212*r - 0.523*g + 0.311*b

	return y, i, q
}

// YIQToRGB converts YIQ to RGB color space.
// Uses NTSC standard for conversion.
// Y should be in range [0,1], I and Q in range [-0.5,0.5].
// Returns RGB values in range [0,1].
func YIQToRGB(y, i, q float64) (r, g, b float64) {
	// NTSC YIQ to RGB conversion matrix
	r = y + 0.9563*i + 0.6210*q
	g = y - 0.2721*i - 0.6474*q
	b = y - 1.1070*i + 1.7046*q

	// Clamp values to [0,1]
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return r, g, b
}
