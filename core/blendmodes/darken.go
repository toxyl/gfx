package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Darken blend mode implementation
type Darken struct {
	BlendMode
}

// NewDarken creates a new Darken blend mode
func NewDarken() *Darken {
	return &Darken{
		BlendMode: BlendMode{
			name: "darken",
		},
	}
}

// Blend implements the darken blend mode
func (d *Darken) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform darken blend: min(a, b)
	r := min(baseLinear.R, blendLinear.R)
	g := min(baseLinear.G, blendLinear.G)
	b := min(baseLinear.B, blendLinear.B)
	a := baseLinear.A

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
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Register the blend mode
func init() {
	Register("darken", "Selects the darker of the base and blend colors for each channel", constants.CategoryDarken, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform darken blend: min(a, b)
		r := min(bottomLinear.R, topLinear.R)
		g := min(bottomLinear.G, topLinear.G)
		b := min(bottomLinear.B, topLinear.B)
		a := bottomLinear.A

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
