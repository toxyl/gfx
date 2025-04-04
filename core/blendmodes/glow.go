package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Glow blend mode implementation
type Glow struct {
	BlendMode
}

// NewGlow creates a new Glow blend mode
func NewGlow() *Glow {
	return &Glow{
		BlendMode: BlendMode{
			name: "glow",
		},
	}
}

// Blend implements the glow blend mode
func (gl *Glow) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform glow blend
	var (
		r     = glowBlend(baseLinear.R, blendLinear.R)
		g     = glowBlend(baseLinear.G, blendLinear.G)
		b     = glowBlend(baseLinear.B, blendLinear.B)
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

// glowBlend performs the glow blend operation for a single channel
func glowBlend(base, blend float64) float64 {
	if base == 0 {
		return 0
	}
	if base == 1 {
		return 1
	}
	// Glow is similar to Reflect but with base and blend colors swapped
	result := blend * blend / (1 - base)
	return math.Min(1.0, result)
}

// Register the blend mode
func init() {
	Register("glow", "Similar to Reflect blend mode but with base and blend colors swapped", constants.CategorySpecial, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform glow blend
		var (
			r = glowBlend(bottomLinear.R, topLinear.R)
			g = glowBlend(bottomLinear.G, topLinear.G)
			b = glowBlend(bottomLinear.B, topLinear.B)
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
