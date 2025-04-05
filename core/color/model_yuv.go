// core/color/base_yuv.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YUV)(nil) // Ensure YUV implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewYUV creates a new YUV instance.
// Y (luminance) is in range [0,1]
// U (blue projection) is in range [-0.436,0.436]
// V (red projection) is in range [-0.615,0.615]
// Alpha is in range [0,1]
func NewYUV[N math.Number](y, u, v, alpha N) (*YUV, error) {
	return newColor(func() *YUV { return &YUV{} }, y, u, v, alpha)
}

// YUVFromRGB converts an RGBA64 (RGB) to a YUV color.
func YUVFromRGB(c *RGBA64) *YUV {
	// Convert RGB to YUV
	y, u, v := utils.RGBToYUV(c.R, c.G, c.B)

	return &YUV{
		Y:     y,
		U:     u,
		V:     v,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// YUV is a helper struct representing a color in the YUV color model.
type YUV struct {
	Y     float64 // [0,1]
	U     float64 // [-0.436,0.436]
	V     float64 // [-0.615,0.615]
	Alpha float64 // [0,1]
}

func (y *YUV) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YUV",
		"Luminance, Blue projection, Red projection color model (BT.601).",
		NewChannelMeta("Y", 0, 1, "", "Luminance."),
		NewChannelMeta("U", -0.436, 0.436, "", "Blue projection."),
		NewChannelMeta("V", -0.615, 0.615, "", "Red projection."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (y *YUV) ToRGB() *RGBA64 {
	// Convert YUV to RGB
	r, g, b := utils.YUVToRGB(y.Y, y.U, y.V)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: y.Alpha,
	}
}

// FromSlice initializes a YUV instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Y, U, V, Alpha.
func (y *YUV) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("YUV requires 4 values, got %d", len(vals))
	}

	y.Y = vals[0]
	y.U = vals[1]
	y.V = vals[2]
	y.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (y *YUV) FromRGBA64(rgba *RGBA64) iColor {
	return YUVFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (y *YUV) ToRGBA64() *RGBA64 {
	return y.ToRGB()
}
