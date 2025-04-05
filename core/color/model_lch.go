// core/color/base_lch.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LCH)(nil) // Ensure LCH implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLCH creates a new LCH instance.
// L (lightness) is in range [0,100]
// C (chroma) is in range [0,100]
// H (hue) is in range [0,360]
// Alpha is in range [0,1]
func NewLCH[N math.Number](l, c, h, alpha N) (*LCH, error) {
	return newColor(func() *LCH { return &LCH{} }, l, c, h, alpha)
}

// LCHFromRGB converts an RGBA64 (RGB) to an LCH color.
func LCHFromRGB(c *RGBA64) *LCH {
	// Convert RGB to LCH
	l, chroma, hue := utils.RGBToLCH(c.R, c.G, c.B)

	return &LCH{
		L:     l,
		C:     chroma,
		H:     hue,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LCH is a helper struct representing a color in the LCH color model.
type LCH struct {
	L     float64 // [0,100] Lightness
	C     float64 // [0,100] Chroma
	H     float64 // [0,360] Hue
	Alpha float64 // [0,1] Alpha
}

func (l *LCH) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LCH",
		"Lightness, Chroma, Hue color model (polar form of LAB).",
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("C", 0, 100, "", "Chroma (colorfulness)."),
		NewChannelMeta("H", 0, 360, "°", "Hue angle."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LCH) ToRGB() *RGBA64 {
	// Convert LCH to RGB
	r, g, b := utils.LCHToRGB(l.L, l.C, l.H)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes an LCH instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: L, C, H, Alpha.
func (l *LCH) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LCH requires 4 values, got %d", len(vals))
	}

	l.L = vals[0]
	l.C = vals[1]
	l.H = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LCH) FromRGBA64(rgba *RGBA64) iColor {
	return LCHFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LCH) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
