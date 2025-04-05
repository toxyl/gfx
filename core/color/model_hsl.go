// core/color/base_hsl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HSL)(nil) // Ensure HSL implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHSL creates a new HSL instance.
// It accepts channel values in the [0,1] range for H, S, L and A.
func NewHSL(h, s, l, alpha float64) (*HSL, error) {
	hsl := &HSL{
		H:     h,
		S:     s,
		L:     l,
		Alpha: alpha,
	}
	return hsl, nil
}

// HSLFromRGB converts an RGBA64 (RGB) to an HSL color.
func HSLFromRGB(c *RGBA64) *HSL {
	h, s, l := utils.RGBToHSL(c.R, c.G, c.B)
	return &HSL{
		H:     h,
		S:     s,
		L:     l,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HSL is a helper struct representing a color in the HSL color model with an alpha channel.
type HSL struct {
	H, S, L, Alpha float64
}

func (hsl *HSL) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HSL",
		"Hue, Saturation, Lightness color model.",
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("S", 0, 100, "%", "Saturation percentage."),
		NewChannelMeta("L", 0, 100, "%", "Lightness percentage."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (hsl *HSL) ToRGB() *RGBA64 {
	r, g, b := utils.HSLToRGB(hsl.H, hsl.S, hsl.L)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: hsl.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (hsl *HSL) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("HSL requires exactly 4 values: H, S, L, Alpha")
	}

	hsl.H = values[0]
	hsl.S = values[1]
	hsl.L = values[2]
	hsl.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (hsl *HSL) FromRGBA64(rgba *RGBA64) iColor {
	return HSLFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (hsl *HSL) ToRGBA64() *RGBA64 {
	return hsl.ToRGB()
}
