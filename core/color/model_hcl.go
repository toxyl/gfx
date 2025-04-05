// core/color/base_hcl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/gfx/core/color/utils"
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
	l2, a, b := utils.LCHToLAB(l, c, h)

	// Then convert to XYZ
	x, y, z := utils.LABToXYZ(l2, a, b)

	// Finally convert to RGB
	r, g, b = utils.XYZToRGB(x, y, z)
	return r, g, b
}

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHCL creates a new HCL instance.
// H (hue) is in range [0,360]
// C (chroma) is in range [0,100]
// L (lightness) is in range [0,100]
// Alpha is in range [0,1]
func NewHCL[N math.Number](h, c, l, alpha N) (*HCL, error) {
	return newColor(func() *HCL { return &HCL{} }, h, c, l, alpha)
}

// HCLFromRGB converts an RGBA64 (RGB) to an HCL color.
func HCLFromRGB(c *RGBA64) *HCL {
	// Convert RGB to HCL
	h, chroma, lightness := utils.RGBToHCL(c.R, c.G, c.B)

	return &HCL{
		H:     h,
		C:     chroma,
		L:     lightness,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HCL is a helper struct representing a color in the HCL color model.
type HCL struct {
	H     float64 // [0,360] Hue
	C     float64 // [0,100] Chroma
	L     float64 // [0,100] Lightness
	Alpha float64 // [0,1] Alpha
}

func (h *HCL) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HCL",
		"Hue, Chroma, Lightness color model (polar form of LAB).",
		NewChannelMeta("H", 0, 360, "°", "Hue angle."),
		NewChannelMeta("C", 0, 100, "", "Chroma (colorfulness)."),
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (h *HCL) ToRGB() *RGBA64 {
	// Convert HCL to RGB
	r, g, b := utils.HCLToRGB(h.H, h.C, h.L)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: h.Alpha,
	}
}

// FromSlice initializes an HCL instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: H, C, L, Alpha.
func (h *HCL) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("HCL requires 4 values, got %d", len(vals))
	}

	h.H = vals[0]
	h.C = vals[1]
	h.L = vals[2]
	h.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (h *HCL) FromRGBA64(rgba *RGBA64) iColor {
	return HCLFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (h *HCL) ToRGBA64() *RGBA64 {
	return h.ToRGB()
}
