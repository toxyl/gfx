package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Average blend mode implementation
type Average struct {
	BlendMode
}

// NewAverage creates a new Average blend mode
func NewAverage() *Average {
	return &Average{
		BlendMode: BlendMode{
			name: "average",
		},
	}
}

// Blend implements the average blend mode
func (a *Average) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform average blend: (a + b) / 2
	var (
		r     = (baseLinear.R + blendLinear.R) / 2
		g     = (baseLinear.G + blendLinear.G) / 2
		b     = (baseLinear.B + blendLinear.B) / 2
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
	Register("average", "Takes the arithmetic mean of the base and blend colors", constants.CategorySpecial, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform average blend: (a + b) / 2
		var (
			r = (bottomLinear.R + topLinear.R) / 2
			g = (bottomLinear.G + topLinear.G) / 2
			b = (bottomLinear.B + topLinear.B) / 2
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
