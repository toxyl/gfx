// core/color/base_hcl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HCL)(nil) // Ensure HCL implements the ColorModel interface.

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
