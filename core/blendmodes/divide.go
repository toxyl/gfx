package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Divide blend mode implementation
type Divide struct {
	BlendMode
}

// NewDivide creates a new Divide blend mode
func NewDivide() *Divide {
	return &Divide{
		BlendMode: BlendMode{
			name: "divide",
		},
	}
}

// Blend implements the divide blend mode
func (d *Divide) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform divide blend: a / b
	var (
		r     = divideBlend(baseLinear.R, blendLinear.R)
		g     = divideBlend(baseLinear.G, blendLinear.G)
		b     = divideBlend(baseLinear.B, blendLinear.B)
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

// divideBlend performs the divide blend operation for a single channel
func divideBlend(base, blend float64) float64 {
	if blend == 0 {
		return 1.0 // Avoid division by zero
	}
	result := base / blend
	if result > 1.0 {
		return 1.0
	}
	return result
}

// Register the blend mode
func init() {
	Register("divide", "Divides the base color by the blend color", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform divide blend: a / b
		var (
			r = divideBlend(bottomLinear.R, topLinear.R)
			g = divideBlend(bottomLinear.G, topLinear.G)
			b = divideBlend(bottomLinear.B, topLinear.B)
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
