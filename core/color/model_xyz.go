// core/color/base_xyz.go
package color

import (
	"fmt"

	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*XYZ)(nil) // Ensure XYZ implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewXYZ creates a new XYZ instance.
// X, Y, Z are tristimulus values in range [0,1]
// Alpha is in range [0,1]
func NewXYZ[N math.Number](x, y, z, alpha N) (*XYZ, error) {
	return newColor(func() *XYZ { return &XYZ{} }, x, y, z, alpha)
}

// XYZFromRGB converts an RGBA64 (RGB) to an XYZ color.
func XYZFromRGB(c *RGBA64) *XYZ {
	x, y, z := rgbToXyz(c.R, c.G, c.B)
	return &XYZ{
		X:     x,
		Y:     y,
		Z:     z,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// XYZ is a helper struct representing a color in the CIE XYZ color model.
// XYZ is a device-independent color space based on human color perception.
type XYZ struct {
	X, Y, Z, Alpha float64
}

func (x *XYZ) Meta() *ColorModelMeta {
	return NewModelMeta(
		"XYZ",
		"CIE XYZ color model based on human color perception.",
		NewChannelMeta("X", 0, 1, "", "X tristimulus value."),
		NewChannelMeta("Y", 0, 1, "", "Y tristimulus value."),
		NewChannelMeta("Z", 0, 1, "", "Z tristimulus value."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (x *XYZ) ToRGB() *RGBA64 {
	r, g, b := xyzToRgb(x.X, x.Y, x.Z)
	return &RGBA64{R: r, G: g, B: b, A: x.Alpha}
}

// FromSlice initializes an XYZ instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: X, Y, Z, Alpha.
func (x *XYZ) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("expected 4 values for XYZ, got %d", len(vals))
	}

	x.X = vals[0]
	x.Y = vals[1]
	x.Z = vals[2]
	x.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (x *XYZ) FromRGBA64(rgba *RGBA64) iColor {
	return XYZFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (x *XYZ) ToRGBA64() *RGBA64 {
	return x.ToRGB()
}
