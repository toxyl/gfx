package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Subtract blend mode implementation.
// This blend mode subtracts the blend color from the base color, ensuring that
// the result never goes below 0. It works in linear color space for accurate
// subtraction calculations.
type Subtract struct {
	BlendMode
}

// NewSubtract creates a new Subtract blend mode.
func NewSubtract() *Subtract {
	return &Subtract{
		BlendMode: BlendMode{
			name: "subtract",
		},
	}
}

// blendSubtract is a helper function that performs the actual subtract blend operation.
// It takes two colors in linear space and returns the blended result.
func blendSubtract(base, blend *color.RGBA64, alpha float64) *color.RGBA64 {
	// Perform subtract blend: max(0, a - b)
	var (
		r = math.Max(0.0, base.R-blend.R)
		g = math.Max(0.0, base.G-blend.G)
		b = math.Max(0.0, base.B-blend.B)
		a = base.A * alpha // Apply alpha to the result
	)

	// Create and return new RGBA64 color
	return &color.RGBA64{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

// Blend implements the subtract blend mode.
func (s *Subtract) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {})

	// Perform the blend operation
	result := blendSubtract(baseLinear, blendLinear, 1.0)

	// Convert back to sRGB space
	return result.Process(false, false, func(c *color.RGBA64) {})
}

// Register the blend mode
func init() {
	Register("subtract", "Subtracts the blend color from the base color", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {})

		// Perform the blend operation
		result := blendSubtract(bottomLinear, topLinear, alpha)

		// Convert back to sRGB space
		return result.Process(false, false, func(c *color.RGBA64) {})
	})
}
