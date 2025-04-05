// core/color/base_cmy.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*CMY)(nil) // Ensure CMY implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewCMY creates a new CMY instance.
// C, M, Y and Alpha are in range [0,1].
func NewCMY[N math.Number](c, m, y, alpha N) (*CMY, error) {
	return newColor(func() *CMY { return &CMY{} }, c, m, y, alpha)
}

// CMYFromRGB converts an RGBA64 (RGB) to a CMY color.
func CMYFromRGB(c *RGBA64) *CMY {
	// Convert RGB to CMY
	cyan, magenta, yellow := utils.RGBToCMY(c.R, c.G, c.B)

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

// CMY is a helper struct representing a color in the CMY color model.
// CMY is a subtractive color model used in color printing.
type CMY struct {
	C     float64 // [0,1] Cyan component
	M     float64 // [0,1] Magenta component
	Y     float64 // [0,1] Yellow component
	Alpha float64 // [0,1] Alpha
}

func (c *CMY) Meta() *ColorModelMeta {
	return NewModelMeta(
		"CMY",
		"Subtractive CMY color model.",
		NewChannelMeta("C", 0, 1, "", "Cyan component."),
		NewChannelMeta("M", 0, 1, "", "Magenta component."),
		NewChannelMeta("Y", 0, 1, "", "Yellow component."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (c *CMY) ToRGB() *RGBA64 {
	// Convert CMY to RGB
	r, g, b := utils.CMYToRGB(c.C, c.M, c.Y)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: c.Alpha,
	}
}

// FromSlice initializes a CMY instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: C, M, Y, Alpha.
func (c *CMY) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("CMY requires 4 values, got %d", len(vals))
	}

	c.C = vals[0]
	c.M = vals[1]
	c.Y = vals[2]
	c.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (c *CMY) FromRGBA64(rgba *RGBA64) iColor {
	return CMYFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (c *CMY) ToRGBA64() *RGBA64 {
	return c.ToRGB()
}
