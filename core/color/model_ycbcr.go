// core/color/base_ycbcr.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YCbCr)(nil) // Ensure YCbCr implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewYCbCr creates a new YCbCr instance.
// Y (luminance) is in range [0,1]
// Cb (blue difference) is in range [-0.5,0.5]
// Cr (red difference) is in range [-0.5,0.5]
// Alpha is in range [0,1]
func NewYCbCr[N math.Number](y, cb, cr, alpha N) (*YCbCr, error) {
	return newColor(func() *YCbCr { return &YCbCr{} }, y, cb, cr, alpha)
}

// YCbCrFromRGB converts an RGBA64 (RGB) to a YCbCr color.
func YCbCrFromRGB(c *RGBA64) *YCbCr {
	// Convert RGB to YCbCr
	y, cb, cr := utils.RGBToYCbCr(c.R, c.G, c.B)

	return &YCbCr{
		Y:     y,
		Cb:    cb,
		Cr:    cr,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// YCbCr is a helper struct representing a color in the YCbCr color model.
type YCbCr struct {
	Y     float64 // [0,1]
	Cb    float64 // [-0.5,0.5]
	Cr    float64 // [-0.5,0.5]
	Alpha float64 // [0,1]
}

func (y *YCbCr) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YCbCr",
		"Luminance, Blue difference, Red difference color model (BT.601).",
		NewChannelMeta("Y", 0, 1, "", "Luminance."),
		NewChannelMeta("Cb", -0.5, 0.5, "", "Blue difference."),
		NewChannelMeta("Cr", -0.5, 0.5, "", "Red difference."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (y *YCbCr) ToRGB() *RGBA64 {
	// Convert YCbCr to RGB
	r, g, b := utils.YCbCrToRGB(y.Y, y.Cb, y.Cr)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: y.Alpha,
	}
}

// FromSlice initializes a YCbCr instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Y, Cb, Cr, Alpha.
func (y *YCbCr) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("YCbCr requires 4 values, got %d", len(vals))
	}

	y.Y = vals[0]
	y.Cb = vals[1]
	y.Cr = vals[2]
	y.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (y *YCbCr) FromRGBA64(rgba *RGBA64) iColor {
	return YCbCrFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (y *YCbCr) ToRGBA64() *RGBA64 {
	return y.ToRGB()
}
