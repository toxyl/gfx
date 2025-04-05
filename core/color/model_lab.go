// core/color/base_lab.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LAB)(nil) // Ensure LAB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLAB creates a new LAB instance.
// L (lightness) is in range [0,100]
// a and b are in range [-128,128]
// Alpha is in range [0,1]
func NewLAB[N math.Number](l, a, b, alpha N) (*LAB, error) {
	return newColor(func() *LAB { return &LAB{} }, l, a, b, alpha)
}

// LABFromRGB converts an RGBA64 (RGB) to a LAB color.
func LABFromRGB(c *RGBA64) *LAB {
	// Convert RGB to LAB
	l, a, b := utils.RGBToLAB(c.R, c.G, c.B)

	return &LAB{
		L:     l,
		A:     a,
		B:     b,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LAB is a helper struct representing a color in the LAB color model.
type LAB struct {
	L     float64 // [0,100] Lightness
	A     float64 // [-128,128] Green-Red
	B     float64 // [-128,128] Blue-Yellow
	Alpha float64 // [0,1] Alpha
}

func (l *LAB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LAB",
		"Lightness, a, b color model (CIELAB).",
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("A", -128, 128, "", "Green-Red component."),
		NewChannelMeta("B", -128, 128, "", "Blue-Yellow component."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LAB) ToRGB() *RGBA64 {
	// Convert LAB to RGB
	r, g, b := utils.LABToRGB(l.L, l.A, l.B)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes a LAB instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: L, A, B, Alpha.
func (l *LAB) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LAB requires 4 values, got %d", len(vals))
	}

	l.L = vals[0]
	l.A = vals[1]
	l.B = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LAB) FromRGBA64(rgba *RGBA64) iColor {
	return LABFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LAB) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
