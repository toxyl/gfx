// core/color/base_cmy.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToCmy(r, g, b float64) (c, m, y float64) {
	// Convert RGB to CMY
	c = constants.CMY_One - r
	m = constants.CMY_One - g
	y = constants.CMY_One - b
	return c, m, y
}

func cmyToRgb(c, m, y float64) (r, g, b float64) {
	// Convert CMY to RGB
	r = constants.CMY_One - c
	g = constants.CMY_One - m
	b = constants.CMY_One - y
	return r, g, b
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*CMY)(nil) // Ensure CMY implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewCMY creates a new CMY instance.
// It accepts channel values in the [0,1] range for C, M, Y and Alpha.
func NewCMY[N math.Number](c, m, y, alpha N) (*CMY, error) {
	return newColor(func() *CMY { return &CMY{} }, c, m, y, alpha)
}

// CMYFromRGB converts an RGBA64 (RGB) to a CMY color.
func CMYFromRGB(c *RGBA64) *CMY {
	cyan, magenta, yellow := rgbToCmy(c.R, c.G, c.B)
	return &CMY{
		C:     cyan,
		M:     magenta,
		Y:     yellow,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// CMY is a helper struct representing a color in the CMY color model with an alpha channel.
type CMY struct {
	C, M, Y, Alpha float64
}

func (cmy *CMY) Meta() *ColorModelMeta {
	return NewModelMeta(
		"CMY",
		"Cyan, Magenta, Yellow color model (subtractive).",
		NewChannelMeta("C", 0, constants.CMY_PercentageMax, "%", "Cyan percentage."),
		NewChannelMeta("M", 0, constants.CMY_PercentageMax, "%", "Magenta percentage."),
		NewChannelMeta("Y", 0, constants.CMY_PercentageMax, "%", "Yellow percentage."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (cmy *CMY) ToRGB() *RGBA64 {
	r, g, b := cmyToRgb(cmy.C, cmy.M, cmy.Y)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: cmy.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (cmy *CMY) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("CMY requires exactly 4 values: C, M, Y, Alpha")
	}

	cmy.C = values[0]
	cmy.M = values[1]
	cmy.Y = values[2]
	cmy.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (cmy *CMY) FromRGBA64(rgba *RGBA64) iColor {
	return CMYFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (cmy *CMY) ToRGBA64() *RGBA64 {
	return cmy.ToRGB()
}
