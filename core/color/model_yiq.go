// core/color/base_yiq.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YIQ)(nil) // Ensure YIQ implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewYIQ creates a new YIQ instance.
// Y (luminance) is in range [0,1]
// I and Q (chrominance) are in range [-0.5,0.5]
// Alpha is in range [0,1]
func NewYIQ[N math.Number](y, i, q, alpha N) (*YIQ, error) {
	return newColor(func() *YIQ { return &YIQ{} }, y, i, q, alpha)
}

// YIQFromRGB converts an RGBA64 (RGB) to a YIQ color.
func YIQFromRGB(c *RGBA64) *YIQ {
	// Convert RGB to YIQ
	y, i, q := utils.RGBToYIQ(c.R, c.G, c.B)

	return &YIQ{
		Y:     y,
		I:     i,
		Q:     q,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// YIQ is a helper struct representing a color in the YIQ color model.
// YIQ is the color space used by NTSC television.
// Y represents luminance, I and Q represent chrominance.
type YIQ struct {
	Y     float64 // [0,1] Luminance
	I     float64 // [-0.5,0.5] In-phase
	Q     float64 // [-0.5,0.5] Quadrature
	Alpha float64 // [0,1] Alpha
}

func (y *YIQ) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YIQ",
		"NTSC YIQ color model.",
		NewChannelMeta("Y", 0, 1, "", "Luminance component."),
		NewChannelMeta("I", -0.5, 0.5, "", "In-phase component."),
		NewChannelMeta("Q", -0.5, 0.5, "", "Quadrature component."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (y *YIQ) ToRGB() *RGBA64 {
	// Convert YIQ to RGB
	r, g, b := utils.YIQToRGB(y.Y, y.I, y.Q)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: y.Alpha,
	}
}

// FromSlice initializes a YIQ instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Y, I, Q, Alpha.
func (y *YIQ) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("YIQ requires 4 values, got %d", len(vals))
	}

	y.Y = vals[0]
	y.I = vals[1]
	y.Q = vals[2]
	y.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (y *YIQ) FromRGBA64(rgba *RGBA64) iColor {
	return YIQFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (y *YIQ) ToRGBA64() *RGBA64 {
	return y.ToRGB()
}
