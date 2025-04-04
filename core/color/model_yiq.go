// core/color/base_yiq.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YIQ)(nil) // Ensure YIQ implements the ColorModel interface.

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

// rgbToYiq converts RGB to YIQ color space.
// YIQ is the color space used by NTSC television.
// Y represents luminance, I and Q represent chrominance.
func rgbToYiq(r, g, b float64) (y, i, q float64) {
	// NTSC YIQ conversion matrix
	y = constants.YIQ_RGB_Y1*r + constants.YIQ_RGB_Y2*g + constants.YIQ_RGB_Y3*b
	i = constants.YIQ_RGB_I1*r + constants.YIQ_RGB_I2*g + constants.YIQ_RGB_I3*b
	q = constants.YIQ_RGB_Q1*r + constants.YIQ_RGB_Q2*g + constants.YIQ_RGB_Q3*b
	return y, i, q
}

// yiqToRgb converts YIQ to RGB color space.
func yiqToRgb(y, i, q float64) (r, g, b float64) {
	// NTSC YIQ to RGB conversion matrix
	// Using inverse matrix coefficients for round trip consistency
	r = math.Clamp(y+constants.YIQ_RGB_R1*i+constants.YIQ_RGB_R2*q, 0, 1)
	g = math.Clamp(y+constants.YIQ_RGB_G1*i+constants.YIQ_RGB_G2*q, 0, 1)
	b = math.Clamp(y+constants.YIQ_RGB_B1*i+constants.YIQ_RGB_B2*q, 0, 1)
	return r, g, b
}

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
	y, i, q := rgbToYiq(c.R, c.G, c.B)
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
type YIQ struct {
	Y, I, Q, Alpha float64
}

func (yiq *YIQ) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YIQ",
		"YIQ color model.",
		NewChannelMeta("Y", 0, 1, "", "Luminance."),
		NewChannelMeta("I", constants.YIQ_I_Min, constants.YIQ_I_Max, "", "In-phase."),
		NewChannelMeta("Q", constants.YIQ_Q_Min, constants.YIQ_Q_Max, "", "Quadrature."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (yiq *YIQ) ToRGB() *RGBA64 {
	r, g, b := yiqToRgb(yiq.Y, yiq.I, yiq.Q)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: yiq.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (yiq *YIQ) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("YIQ requires exactly 4 values: Y, I, Q, Alpha")
	}

	yiq.Y = values[0]
	yiq.I = values[1]
	yiq.Q = values[2]
	yiq.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (yiq *YIQ) FromRGBA64(rgba *RGBA64) iColor {
	return YIQFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (yiq *YIQ) ToRGBA64() *RGBA64 {
	return yiq.ToRGB()
}
