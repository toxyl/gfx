package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Negation blend mode implementation
type Negation struct {
	BlendMode
}

// NewNegation creates a new Negation blend mode
func NewNegation() *Negation {
	return &Negation{
		BlendMode: BlendMode{
			name: "negation",
		},
	}
}

// Blend implements the negation blend mode
func (n *Negation) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform negation blend
	var (
		r     = negationBlend(baseLinear.R, blendLinear.R)
		g     = negationBlend(baseLinear.G, blendLinear.G)
		b     = negationBlend(baseLinear.B, blendLinear.B)
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

// negationBlend performs the negation blend operation for a single channel
func negationBlend(base, blend float64) float64 {
	// Negation: 1 - abs(1 - base - blend)
	result := 1 - math.Abs(1-base-blend)
	return math.Min(1.0, math.Max(0.0, result))
}

// Register the blend mode
func init() {
	Register("negation", "Inverts the result of the Difference blend mode", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform negation blend
		var (
			r = negationBlend(bottomLinear.R, topLinear.R)
			g = negationBlend(bottomLinear.G, topLinear.G)
			b = negationBlend(bottomLinear.B, topLinear.B)
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
