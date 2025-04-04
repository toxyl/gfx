// core/color/base_hcl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HCL)(nil) // Ensure HCL implements the ColorModel interface.

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

// labToHcl converts LAB to HCL color space.
func labToHcl(l, a, b float64) (h, c, l2 float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.HCL_WhiteThreshold && math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, 100
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.HCL_NearWhiteThreshold && math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, l
	}

	// Handle gray colors (a and b close to zero)
	if math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, l
	}

	c = math.Sqrt(a*a + b*b)

	// Handle very small chroma values
	if c < constants.HCL_ChromaThreshold {
		return 0, 0, l
	}

	h = math.Atan2(b, a) * constants.HCL_DegreesPerRadian

	if h < 0 {
		h += constants.HCL_DegreesPerCircle
	}

	return h, c, l
}

// hclToLab converts HCL to LAB color space.
func hclToLab(h, c, l float64) (l2, a, b float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.HCL_WhiteThreshold && math.Abs(c) < constants.HCL_WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.HCL_NearWhiteThreshold && math.Abs(c) < constants.HCL_WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (chroma close to zero)
	if math.Abs(c) < constants.HCL_ChromaThreshold {
		return l, 0, 0
	}

	h = h * constants.HCL_RadiansPerDegree
	a = c * math.Cos(h)
	b = c * math.Sin(h)

	return l, a, b
}

// hclToRgb converts HCL to RGB color space.
func hclToRgb(h, c, l float64) (r, g, b float64) {
	// First convert to LAB
	l2, a, b := hclToLab(h, c, l)

	// Then convert to XYZ
	x, y, z := labToXyz(l2, a, b)

	// Finally convert to RGB
	r, g, b = xyzToRgb(x, y, z)
	return r, g, b
}

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHCL creates a new HCL instance.
// H (hue) is in range [0,360]
// C (chroma) is in range [0,1]
// L (luminance) is in range [0,1]
// Alpha is in range [0,1]
func NewHCL[N math.Number](h, c, l, alpha N) (*HCL, error) {
	return newColor(func() *HCL { return &HCL{} }, h, c, l, alpha)
}

// HCLFromRGB converts an RGBA64 (RGB) to an HCL color.
func HCLFromRGB(c *RGBA64) *HCL {
	// First convert to LAB
	x, y, z := rgbToXyz(c.R, c.G, c.B)
	l, a, b := xyzToLab(x, y, z)

	// Then convert to HCL
	h, chroma, l2 := labToHcl(l, a, b)

	return &HCL{
		H:     h,
		C:     chroma,
		L:     l2,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HCL is a helper struct representing a color in the HCL color model.
// HCL is a polar version of LAB that's more intuitive for some applications.
type HCL struct {
	H, C, L, Alpha float64
}

func (hcl *HCL) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HCL",
		"Hue, Chroma, Luminance color model.",
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("C", 0, 1, "", "Chroma."),
		NewChannelMeta("L", 0, 1, "", "Luminance."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (hcl *HCL) ToRGB() *RGBA64 {
	r, g, b := hclToRgb(hcl.H, hcl.C, hcl.L)
	return &RGBA64{R: r, G: g, B: b, A: hcl.Alpha}
}

// FromSlice initializes an HCL instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: H, C, L, Alpha.
func (hcl *HCL) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("HCL requires 4 values, got %d", len(vals))
	}

	hcl.H = vals[0]
	hcl.C = vals[1]
	hcl.L = vals[2]
	hcl.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (hcl *HCL) FromRGBA64(rgba *RGBA64) iColor {
	return HCLFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (hcl *HCL) ToRGBA64() *RGBA64 {
	return hcl.ToRGB()
}
