package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LUV)(nil) // Ensure LUV implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLUV creates a new LUV instance.
// L (lightness) is in range [0,100]
// u and v are in range [-100,100]
// Alpha is in range [0,1]
func NewLUV[N math.Number](l, u, v, alpha N) (*LUV, error) {
	return newColor(func() *LUV { return &LUV{} }, l, u, v, alpha)
}

// LUVFromRGB converts an RGBA64 (RGB) to an LUV color.
func LUVFromRGB(c *RGBA64) *LUV {
	// Convert RGB to LUV
	l, u, v := utils.RGBToLUV(c.R, c.G, c.B)

	return &LUV{
		L:     l,
		U:     u,
		V:     v,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LUV is a helper struct representing a color in the LUV color model.
type LUV struct {
	L     float64 // [0,100] Lightness
	U     float64 // [-100,100] u component
	V     float64 // [-100,100] v component
	Alpha float64 // [0,1] Alpha
}

func (l *LUV) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LUV",
		"Lightness, u, v color model (CIELUV).",
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("U", -100, 100, "", "u component."),
		NewChannelMeta("V", -100, 100, "", "v component."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LUV) ToRGB() *RGBA64 {
	// Convert LUV to RGB
	r, g, b := utils.LUVToRGB(l.L, l.U, l.V)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes an LUV instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: L, U, V, Alpha.
func (l *LUV) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LUV requires 4 values, got %d", len(vals))
	}

	l.L = vals[0]
	l.U = vals[1]
	l.V = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LUV) FromRGBA64(rgba *RGBA64) iColor {
	return LUVFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LUV) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
