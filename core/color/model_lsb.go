// core/color/base_lambdasb.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LSB)(nil) // Ensure LSB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLSB creates a new LSB instance.
// Lightness is in range [0,1]
// Saturation is in range [0,1]
// Brightness is in range [0,1]
// Alpha is in range [0,1]
func NewLSB[N math.Number](lightness, saturation, brightness, alpha N) (*LSB, error) {
	return newColor(func() *LSB { return &LSB{} }, lightness, saturation, brightness, alpha)
}

// LSBFromRGB converts an RGBA64 (RGB) to an LSB color.
func LSBFromRGB(c *RGBA64) *LSB {
	// Convert RGB to HSL first
	_, s, l := utils.RGBToHSL(c.R, c.G, c.B)

	// Convert RGB to HSB to get brightness
	_, _, b := utils.RGBToHSB(c.R, c.G, c.B)

	return &LSB{
		Lightness:  l,
		Saturation: s,
		Brightness: b,
		Alpha:      c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LSB is a helper struct representing a color in the LSB color model.
type LSB struct {
	Lightness  float64 // [0,1]
	Saturation float64 // [0,1]
	Brightness float64 // [0,1]
	Alpha      float64 // [0,1]
}

func (l *LSB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LSB",
		"Lightness, Saturation, Brightness color model.",
		NewChannelMeta("L", 0, 1, "", "Lightness."),
		NewChannelMeta("S", 0, 1, "", "Saturation."),
		NewChannelMeta("B", 0, 1, "", "Brightness."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LSB) ToRGB() *RGBA64 {
	// Convert LSB to RGB via HSL
	r, g, b := utils.HSLToRGB(0, l.Saturation, l.Lightness)

	// Adjust brightness
	if l.Brightness < 1 {
		r *= l.Brightness
		g *= l.Brightness
		b *= l.Brightness
	}

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes a LSB instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Lightness, Saturation, Brightness, Alpha.
func (l *LSB) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LSB requires 4 values, got %d", len(vals))
	}

	l.Lightness = vals[0]
	l.Saturation = vals[1]
	l.Brightness = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LSB) FromRGBA64(rgba *RGBA64) iColor {
	return LSBFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LSB) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
