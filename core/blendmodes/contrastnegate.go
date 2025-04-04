package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// ContrastNegate blend mode implementation
type ContrastNegate struct {
	BlendMode
}

// NewContrastNegate creates a new ContrastNegate blend mode
func NewContrastNegate() *ContrastNegate {
	return &ContrastNegate{
		BlendMode: BlendMode{
			name: "contrastnegate",
		},
	}
}

// Blend implements the contrast-negate blend mode
func (cn *ContrastNegate) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform contrast-negate blend
	var (
		r     = contrastNegateBlend(baseLinear.R, blendLinear.R)
		g     = contrastNegateBlend(baseLinear.G, blendLinear.G)
		b     = contrastNegateBlend(baseLinear.B, blendLinear.B)
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

// contrastNegateBlend performs the contrast-negate blend operation for a single channel
func contrastNegateBlend(base, blend float64) float64 {
	// ContrastNegate: abs(base - blend) + 0.5
	result := math.Abs(base-blend) + 0.5
	return math.Min(1.0, result)
}

// Register the blend mode
func init() {
	Register("contrastnegate", "Creates contrast by taking the absolute difference between base and blend colors and adding 0.5", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform contrast-negate blend
		var (
			r = contrastNegateBlend(bottomLinear.R, topLinear.R)
			g = contrastNegateBlend(bottomLinear.G, topLinear.G)
			b = contrastNegateBlend(bottomLinear.B, topLinear.B)
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
