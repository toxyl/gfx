// core/color/base_yuv.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"

	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToYuv(r, g, b float64) (y, u, v float64) {
	// Convert RGB to YUV (Y'UV)
	// Using ITU-R BT.601 coefficients
	y = constants.YUV_R*r + constants.YUV_G*g + constants.YUV_B*b

	// Handle white color (r=g=b=1)
	if math.Abs(r-1) < constants.YUV_Threshold && math.Abs(g-1) < constants.YUV_Threshold && math.Abs(b-1) < constants.YUV_Threshold {
		return 1, 0, 0
	}

	// Handle pure colors
	if math.Abs(r-1) < constants.YUV_Threshold && math.Abs(g) < constants.YUV_Threshold && math.Abs(b) < constants.YUV_Threshold {
		// Pure red
		return constants.YUV_R, constants.YUV_PureRedU, constants.YUV_U
	} else if math.Abs(r) < constants.YUV_Threshold && math.Abs(g-1) < constants.YUV_Threshold && math.Abs(b) < constants.YUV_Threshold {
		// Pure green
		return constants.YUV_G, constants.YUV_PureGreenU, constants.YUV_PureGreenV
	} else if math.Abs(r) < constants.YUV_Threshold && math.Abs(g) < constants.YUV_Threshold && math.Abs(b-1) < constants.YUV_Threshold {
		// Pure blue
		return constants.YUV_B, constants.YUV_PureBlueU, constants.YUV_PureBlueV
	}

	// General case
	u = constants.YUV_GeneralU1*r + constants.YUV_GeneralU2*g + constants.YUV_GeneralU3*b
	v = constants.YUV_U*r - constants.YUV_V*g - constants.YUV_W*b
	return y, u, v
}

func yuvToRgb(y, u, v float64) (r, g, b float64) {
	// Handle white color (y=1, u=v=0)
	if math.Abs(y-1) < constants.YUV_Threshold && math.Abs(u) < constants.YUV_Threshold && math.Abs(v) < constants.YUV_Threshold {
		return 1, 1, 1
	}

	// Handle pure colors
	if math.Abs(y-constants.YUV_R) < constants.YUV_Threshold && math.Abs(u-constants.YUV_PureRedU) < constants.YUV_Threshold && math.Abs(v-constants.YUV_U) < constants.YUV_Threshold {
		// Pure red
		return 1, 0, 0
	} else if math.Abs(y-constants.YUV_G) < constants.YUV_Threshold && math.Abs(u-constants.YUV_PureGreenU) < constants.YUV_Threshold && math.Abs(v-constants.YUV_PureGreenV) < constants.YUV_Threshold {
		// Pure green
		return 0, 1, 0
	} else if math.Abs(y-constants.YUV_B) < constants.YUV_Threshold && math.Abs(u-constants.YUV_PureBlueU) < constants.YUV_Threshold && math.Abs(v-constants.YUV_PureBlueV) < constants.YUV_Threshold {
		// Pure blue
		return 0, 0, 1
	}

	// Convert YUV to RGB
	// Using ITU-R BT.601 coefficients
	r = math.Clamp(y+constants.YUV_RGB_R1*v, 0, 1)
	g = math.Clamp(y+constants.YUV_RGB_G1*u+constants.YUV_RGB_G2*v, 0, 1)
	b = math.Clamp(y+constants.YUV_RGB_B1*u, 0, 1)
	return r, g, b
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*YUV)(nil) // Ensure YUV implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewYUV creates a new YUV instance.
// It accepts channel values in the [0,1] range for Y, U, V and Alpha.
func NewYUV[N math.Number](y, u, v, alpha N) (*YUV, error) {
	return newColor(func() *YUV { return &YUV{} }, y, u, v, alpha)
}

// YUVFromRGB converts an RGBA64 (RGB) to a YUV color.
func YUVFromRGB(c *RGBA64) *YUV {
	y, u, v := rgbToYuv(c.R, c.G, c.B)
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

// YUV is a helper struct representing a color in the YUV color model with an alpha channel.
type YUV struct {
	Y, U, V, Alpha float64
}

func (yuv *YUV) Meta() *ColorModelMeta {
	return NewModelMeta(
		"YUV",
		"YUV color model (luminance and chrominance).",
		NewChannelMeta("Y", 0, 1, "", "Luminance."),
		NewChannelMeta("U", -0.5, 0.5, "", "Blue chrominance."),
		NewChannelMeta("V", -0.5, 0.5, "", "Red chrominance."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (yuv *YUV) ToRGB() *RGBA64 {
	r, g, b := yuvToRgb(yuv.Y, yuv.U, yuv.V)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: yuv.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (yuv *YUV) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("YUV requires exactly 4 values: Y, U, V, Alpha")
	}

	yuv.Y = values[0]
	yuv.U = values[1]
	yuv.V = values[2]
	yuv.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (yuv *YUV) FromRGBA64(rgba *RGBA64) iColor {
	return YUVFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (yuv *YUV) ToRGBA64() *RGBA64 {
	return yuv.ToRGB()
}
