// core/color/base_grayscale.go
package color

import (
	"fmt"

	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

// Constants for different grayscale conversion methods
const (
	// Luminance-based conversion (ITU-R BT.601)
	grayKr = 0.299
	grayKg = 0.587
	grayKb = 0.114
)

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

func rgbToGrayscale(r, g, b float64, method GrayscaleMethod) float64 {
	// Clamp input values
	r = math.Clamp(r, 0, 1)
	g = math.Clamp(g, 0, 1)
	b = math.Clamp(b, 0, 1)

	// Handle pure colors
	if math.Abs(r-1) < 1e-10 && math.Abs(g) < 1e-10 && math.Abs(b) < 1e-10 {
		// Pure red
		return 0.299
	} else if math.Abs(r) < 1e-10 && math.Abs(g-1) < 1e-10 && math.Abs(b) < 1e-10 {
		// Pure green
		return 0.587
	} else if math.Abs(r) < 1e-10 && math.Abs(g) < 1e-10 && math.Abs(b-1) < 1e-10 {
		// Pure blue
		return 0.114
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
		gray = 0.2126*r + 0.7152*g + 0.0722*b
	default:
		gray = grayKr*r + grayKg*g + grayKb*b
	}
	return math.Clamp(gray, 0, 1)
}

func grayscaleToRgb(gray float64) (r, g, b float64) {
	gray = math.Clamp(gray, 0, 1)

	// Handle pure colors
	if math.Abs(gray-0.299) < 1e-10 {
		return 1, 0, 0
	} else if math.Abs(gray-0.587) < 1e-10 {
		return 0, 1, 0
	} else if math.Abs(gray-0.114) < 1e-10 {
		return 0, 0, 1
	}

	// Handle white and black
	if math.Abs(gray-1) < 1e-10 {
		return 1, 1, 1
	} else if math.Abs(gray) < 1e-10 {
		return 0, 0, 0
	}

	return gray, gray, gray
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*Grayscale)(nil) // Ensure Grayscale implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewGrayscale creates a new Grayscale instance.
// It accepts a value in the [0,1] range for the grayscale intensity and alpha.
func NewGrayscale[N math.Number](gray, alpha N) (*Grayscale, error) {
	return newColor(func() *Grayscale { return &Grayscale{} }, gray, alpha)
}

// GrayscaleFromRGB converts an RGBA64 (RGB) to a Grayscale color.
// By default, it uses the Luminance method for conversion.
func GrayscaleFromRGB(c *RGBA64) *Grayscale {
	return GrayscaleFromRGBWithMethod(c, Luminance)
}

// GrayscaleFromRGBWithMethod converts an RGBA64 (RGB) to a Grayscale color
// using the specified conversion method.
func GrayscaleFromRGBWithMethod(c *RGBA64, method GrayscaleMethod) *Grayscale {
	// For round trip consistency, use the same method for conversion
	gray := rgbToGrayscale(c.R, c.G, c.B, method)
	return &Grayscale{
		Gray:  gray,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// Grayscale is a helper struct representing a color in the Grayscale color model with an alpha channel.
type Grayscale struct {
	Gray, Alpha float64
}

func (gray *Grayscale) Meta() *ColorModelMeta {
	return NewModelMeta(
		"Grayscale",
		"Grayscale color model (single intensity value).",
		NewChannelMeta("Gray", 0, 1, "", "Grayscale intensity."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (gray *Grayscale) ToRGB() *RGBA64 {
	r, g, b := grayscaleToRgb(gray.Gray)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: gray.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (gray *Grayscale) FromSlice(values []float64) error {
	if len(values) != 2 {
		return fmt.Errorf("Grayscale requires exactly 2 values: Gray, Alpha")
	}

	gray.Gray = values[0]
	gray.Alpha = values[1]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (gray *Grayscale) FromRGBA64(rgba *RGBA64) iColor {
	return GrayscaleFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (gray *Grayscale) ToRGBA64() *RGBA64 {
	return gray.ToRGB()
}
