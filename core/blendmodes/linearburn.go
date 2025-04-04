package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// LinearBurn blend mode implementation
type LinearBurn struct {
	BlendMode
}

// NewLinearBurn creates a new LinearBurn blend mode
func NewLinearBurn() *LinearBurn {
	return &LinearBurn{
		BlendMode: BlendMode{
			name: "linearburn",
		},
	}
}

// Blend implements the linear burn blend mode
func (lb *LinearBurn) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform linear burn blend: a + b - 1
	var (
		r     = math.Max(0.0, baseLinear.R+blendLinear.R-1.0)
		g     = math.Max(0.0, baseLinear.G+blendLinear.G-1.0)
		b     = math.Max(0.0, baseLinear.B+blendLinear.B-1.0)
		alpha = baseLinear.A
	)

	// Create and return new RGBA64 color
	result := &color.RGBA64{
		R: r,
		G: g,
		B: b,
		A: alpha,
	}

	// Convert back to sRGB space
	return result.Process(false, false, func(c *color.RGBA64) {
		// No additional processing needed
	})
}

// Register the blend mode
func init() {
	Register("linearburn", "Darkens the base color by decreasing brightness", constants.CategoryDarken, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform linear burn blend: a + b - 1
		var (
			r = math.Max(0.0, bottomLinear.R+topLinear.R-1.0)
			g = math.Max(0.0, bottomLinear.G+topLinear.G-1.0)
			b = math.Max(0.0, bottomLinear.B+topLinear.B-1.0)
			a = bottomLinear.A
		)

		// Create and return new RGBA64 color
		result := &color.RGBA64{
			R: r,
			G: g,
			B: b,
			A: a,
		}

		// Convert back to sRGB space
		return result.Process(false, false, func(c *color.RGBA64) {
			// No additional processing needed
		})
	})
}
