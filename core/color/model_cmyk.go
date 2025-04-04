// core/color/base_cmyk.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToCmyk(r, g, b float64) (c, m, y, k float64) {
	k = constants.CMYK_One - math.Max(r, math.Max(g, b))
	if k == constants.CMYK_One {
		return 0, 0, 0, constants.CMYK_One
	}
	c = (constants.CMYK_One - r - k) / (constants.CMYK_One - k)
	m = (constants.CMYK_One - g - k) / (constants.CMYK_One - k)
	y = (constants.CMYK_One - b - k) / (constants.CMYK_One - k)
	return
}

func cmykToRgb(c, m, y, k float64) (r, g, b float64) {
	r = (constants.CMYK_One - c) * (constants.CMYK_One - k)
	g = (constants.CMYK_One - m) * (constants.CMYK_One - k)
	b = (constants.CMYK_One - y) * (constants.CMYK_One - k)
	return
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*CMYK)(nil) // Ensure CMYK implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewCMYK creates a new CMYK instance.
// It accepts channel values in the [0,1] range for C, M, Y and A.
func NewCMYK[N math.Number](c, m, y, k, a N) (*CMYK, error) {
	return newColor(func() *CMYK { return &CMYK{} }, c, m, y, k, a)
}

// CMYKFromRGB converts an RGBA64 (RGB) to a CMYK color.
func CMYKFromRGB(c *RGBA64) *CMYK {
	cVal, mVal, yVal, kVal := rgbToCmyk(c.R, c.G, c.B)
	return &CMYK{
		C: cVal,
		M: mVal,
		Y: yVal,
		K: kVal,
		A: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// CMYK is a helper struct representing a color in the CMYK color model with an alpha channel.
type CMYK struct {
	C, M, Y, K, A float64
}

func (cmyk *CMYK) Meta() *ColorModelMeta {
	return NewModelMeta(
		"CMYK",
		"Cyan, Magenta, Yellow, Black color model.",
		NewChannelMeta("C", 0, constants.CMYK_PercentageMax, "%", "Cyan percentage."),
		NewChannelMeta("M", 0, constants.CMYK_PercentageMax, "%", "Magenta percentage."),
		NewChannelMeta("Y", 0, constants.CMYK_PercentageMax, "%", "Yellow percentage."),
		NewChannelMeta("K", 0, constants.CMYK_PercentageMax, "%", "Key (black) percentage."),
		NewChannelMeta("A", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (cmyk *CMYK) ToRGB() *RGBA64 {
	r, g, b := cmykToRgb(cmyk.C, cmyk.M, cmyk.Y, cmyk.K)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: cmyk.A,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (cmyk *CMYK) FromSlice(values []float64) error {
	if len(values) != 5 {
		return fmt.Errorf("CMYK requires exactly 5 values: C, M, Y, K, Alpha")
	}

	cmyk.C = values[0]
	cmyk.M = values[1]
	cmyk.Y = values[2]
	cmyk.K = values[3]
	cmyk.A = values[4]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (cmyk *CMYK) FromRGBA64(rgba *RGBA64) iColor {
	return CMYKFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (cmyk *CMYK) ToRGBA64() *RGBA64 {
	return cmyk.ToRGB()
}
