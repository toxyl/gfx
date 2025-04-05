// core/color/base_grayscale.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*Grayscale)(nil) // Ensure Grayscale implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewGrayscale creates a new Grayscale instance.
// Gray and Alpha are in range [0,1].
func NewGrayscale[N math.Number](gray, alpha N) (*Grayscale, error) {
	return newColor(func() *Grayscale { return &Grayscale{} }, gray, alpha)
}

// GrayscaleFromRGB converts an RGBA64 (RGB) to a Grayscale color.
// Uses the Luminance method (ITU-R BT.601) by default.
func GrayscaleFromRGB(c *RGBA64) *Grayscale {
	// Convert RGB to Grayscale using Luminance method
	gray := utils.RGBToGrayscale(c.R, c.G, c.B, utils.Luminance)

	return &Grayscale{
		Gray:  gray,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// Grayscale is a helper struct representing a color in the Grayscale color model.
// Grayscale represents colors using a single intensity value.
type Grayscale struct {
	Gray  float64 // [0,1] Grayscale intensity
	Alpha float64 // [0,1] Alpha
}

func (g *Grayscale) Meta() *ColorModelMeta {
	return NewModelMeta(
		"Grayscale",
		"Single-channel grayscale color model.",
		NewChannelMeta("Gray", 0, 1, "", "Grayscale intensity."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (g *Grayscale) ToRGB() *RGBA64 {
	// Convert Grayscale to RGB
	r, gVal, b := utils.GrayscaleToRGB(g.Gray)

	return &RGBA64{
		R: r,
		G: gVal,
		B: b,
		A: g.Alpha,
	}
}

// FromSlice initializes a Grayscale instance from a slice of float64 values.
// The slice must contain exactly 2 values in the order: Gray, Alpha.
func (g *Grayscale) FromSlice(vals []float64) error {
	if len(vals) != 2 {
		return fmt.Errorf("Grayscale requires 2 values, got %d", len(vals))
	}

	g.Gray = vals[0]
	g.Alpha = vals[1]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (g *Grayscale) FromRGBA64(rgba *RGBA64) iColor {
	return GrayscaleFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (g *Grayscale) ToRGBA64() *RGBA64 {
	return g.ToRGB()
}
