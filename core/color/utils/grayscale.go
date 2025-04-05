package utils

import "github.com/toxyl/math"

// GrayscaleMethod defines the method used to convert RGB to grayscale
type GrayscaleMethod int

const (
	// Luminance uses the ITU-R BT.601 standard weights
	Luminance GrayscaleMethod = iota
	// Average uses equal weights for all channels
	Average
	// Lightness uses the HSL lightness value
	Lightness
	// Luminosity uses the ITU-R BT.709 standard weights
	Luminosity
)

// Constants for different grayscale conversion methods
const (
	// Luminance-based conversion (ITU-R BT.601)
	grayKr = 0.299
	grayKg = 0.587
	grayKb = 0.114

	// Luminosity-based conversion (ITU-R BT.709)
	grayLr = 0.2126
	grayLg = 0.7152
	grayLb = 0.0722
)

// RGBToGrayscale converts RGB to grayscale using the specified method.
// RGB values should be in range [0,1].
// Returns grayscale value in range [0,1].
func RGBToGrayscale(r, g, b float64, method GrayscaleMethod) float64 {
	// Clamp input values
	r = math.Clamp(r, 0, 1)
	g = math.Clamp(g, 0, 1)
	b = math.Clamp(b, 0, 1)

	// Handle pure colors
	if math.Abs(r-1) < 1e-10 && math.Abs(g) < 1e-10 && math.Abs(b) < 1e-10 {
		// Pure red
		return grayKr
	} else if math.Abs(r) < 1e-10 && math.Abs(g-1) < 1e-10 && math.Abs(b) < 1e-10 {
		// Pure green
		return grayKg
	} else if math.Abs(r) < 1e-10 && math.Abs(g) < 1e-10 && math.Abs(b-1) < 1e-10 {
		// Pure blue
		return grayKb
	}

	// Handle white
	if math.Abs(r-1) < 1e-10 && math.Abs(g-1) < 1e-10 && math.Abs(b-1) < 1e-10 {
		return 1
	}

	// Handle black
	if math.Abs(r) < 1e-10 && math.Abs(g) < 1e-10 && math.Abs(b) < 1e-10 {
		return 0
	}

	var gray float64
	switch method {
	case Luminance:
		gray = grayKr*r + grayKg*g + grayKb*b
	case Average:
		gray = (r + g + b) / 3
	case Lightness:
		max := math.Max(math.Max(r, g), b)
		min := math.Min(math.Min(r, g), b)
		gray = (max + min) / 2
	case Luminosity:
		gray = grayLr*r + grayLg*g + grayLb*b
	default:
		gray = grayKr*r + grayKg*g + grayKb*b
	}
	return math.Clamp(gray, 0, 1)
}

// GrayscaleToRGB converts a grayscale value to RGB.
// Gray value should be in range [0,1].
// Returns RGB values in range [0,1].
func GrayscaleToRGB(gray float64) (r, g, b float64) {
	// Clamp input value
	gray = math.Clamp(gray, 0, 1)
	return gray, gray, gray
}
