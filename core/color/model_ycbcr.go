// core/color/base_ycbcr.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToYCbCr(r, g, b float64) (y, cb, cr float64) {
	// Convert RGB to YCbCr
	y = constants.YCBCR_KR*r + constants.YCBCR_KG*g + constants.YCBCR_KB*b
	cb = (b - y) / (2 * (1 - constants.YCBCR_KB))
	cr = (r - y) / (2 * (1 - constants.YCBCR_KR))
	return y, cb, cr
}

func ycbcrToRgb(y, cb, cr float64) (r, g, b float64) {
	// Convert YCbCr to RGB
	r = y + cr*(1-constants.YCBCR_KR)/constants.YCBCR_Divisor
	b = y + cb*(1-constants.YCBCR_KB)/constants.YCBCR_Divisor
	g = (y - constants.YCBCR_KR*r - constants.YCBCR_KB*b) / constants.YCBCR_KG
	return r, g, b
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YCbCr)(nil) // Ensure YCbCr implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewYCbCr creates a new YCbCr instance.
// It accepts channel values in the [0,1] range for Y, Cb, Cr and Alpha.
func NewYCbCr[N math.Number](y, cb, cr, alpha N) (*YCbCr, error) {
	return newColor(func() *YCbCr { return &YCbCr{} }, y, cb, cr, alpha)
}

// YCbCrFromRGB converts an RGBA64 (RGB) to a YCbCr color.
func YCbCrFromRGB(c *RGBA64) *YCbCr {
	y, cb, cr := rgbToYCbCr(c.R, c.G, c.B)
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

// YCbCr is a helper struct representing a color in the YCbCr color model with an alpha channel.
type YCbCr struct {
	Y, Cb, Cr, Alpha float64
}

func (ycbcr *YCbCr) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YCbCr",
		"YCbCr color model (digital YUV).",
		NewChannelMeta("Y", 0, 1, "", "Luminance."),
		NewChannelMeta("Cb", constants.YCBCR_Cb_Min, constants.YCBCR_Cb_Max, "", "Blue chrominance."),
		NewChannelMeta("Cr", constants.YCBCR_Cr_Min, constants.YCBCR_Cr_Max, "", "Red chrominance."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (ycbcr *YCbCr) ToRGB() *RGBA64 {
	r, g, b := ycbcrToRgb(ycbcr.Y, ycbcr.Cb, ycbcr.Cr)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: ycbcr.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (ycbcr *YCbCr) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("YCbCr requires exactly 4 values: Y, Cb, Cr, Alpha")
	}

	ycbcr.Y = values[0]
	ycbcr.Cb = values[1]
	ycbcr.Cr = values[2]
	ycbcr.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (ycbcr *YCbCr) FromRGBA64(rgba *RGBA64) iColor {
	return YCbCrFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (ycbcr *YCbCr) ToRGBA64() *RGBA64 {
	return ycbcr.ToRGB()
}
