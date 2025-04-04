package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Erase blend mode implementation
type Erase struct {
	BlendMode
}

// NewErase creates a new Erase blend mode
func NewErase() *Erase {
	return &Erase{
		BlendMode: BlendMode{
			name: "erase",
		},
	}
}

// Blend implements the erase blend mode
func (e *Erase) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform erase blend
	var (
		r     = baseLinear.R
		g     = baseLinear.G
		b     = baseLinear.B
		alpha = baseLinear.A * (1 - blendLinear.A)
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
	Register("erase", "Erases the base color based on the blend color's alpha channel", constants.CategoryBasic, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform erase blend
		var (
			r = bottomLinear.R
			g = bottomLinear.G
			b = bottomLinear.B
			a = bottomLinear.A * (1 - topLinear.A)
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
