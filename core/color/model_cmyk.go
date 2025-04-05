// core/color/base_cmyk.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*CMYK)(nil) // Ensure CMYK implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewCMYK creates a new CMYK instance.
// All components (Cyan, Magenta, Yellow, Key) and Alpha are in range [0,1].
func NewCMYK[N math.Number](c, m, y, k, alpha N) (*CMYK, error) {
	return newColor(func() *CMYK { return &CMYK{} }, c, m, y, k, alpha)
}

// CMYKFromRGB converts an RGBA64 (RGB) to a CMYK color.
func CMYKFromRGB(c *RGBA64) *CMYK {
	// Convert RGB to CMYK
	cyan, magenta, yellow, key := utils.RGBToCMYK(c.R, c.G, c.B)

	return &CMYK{
		C:     cyan,
		M:     magenta,
		Y:     yellow,
		K:     key,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// CMYK is a helper struct representing a color in the CMYK color model.
type CMYK struct {
	C     float64 // [0,1] Cyan
	M     float64 // [0,1] Magenta
	Y     float64 // [0,1] Yellow
	K     float64 // [0,1] Key (Black)
	Alpha float64 // [0,1] Alpha
}

func (c *CMYK) Meta() *ColorModelMeta {
	return NewModelMeta(
		"CMYK",
		"Cyan, Magenta, Yellow, Key color model (subtractive color model).",
		NewChannelMeta("C", 0, 1, "", "Cyan component."),
		NewChannelMeta("M", 0, 1, "", "Magenta component."),
		NewChannelMeta("Y", 0, 1, "", "Yellow component."),
		NewChannelMeta("K", 0, 1, "", "Key (Black) component."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (c *CMYK) ToRGB() *RGBA64 {
	// Convert CMYK to RGB
	r, g, b := utils.CMYKToRGB(c.C, c.M, c.Y, c.K)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: c.Alpha,
	}
}

// FromSlice initializes a CMYK instance from a slice of float64 values.
// The slice must contain exactly 5 values in the order: C, M, Y, K, Alpha.
func (c *CMYK) FromSlice(vals []float64) error {
	if len(vals) != 5 {
		return fmt.Errorf("CMYK requires 5 values, got %d", len(vals))
	}

	c.C = vals[0]
	c.M = vals[1]
	c.Y = vals[2]
	c.K = vals[3]
	c.Alpha = vals[4]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (c *CMYK) FromRGBA64(rgba *RGBA64) iColor {
	return CMYKFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (c *CMYK) ToRGBA64() *RGBA64 {
	return c.ToRGB()
}
