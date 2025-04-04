package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Exclusion blend mode implementation
type Exclusion struct {
	BlendMode
}

// NewExclusion creates a new Exclusion blend mode
func NewExclusion() *Exclusion {
	return &Exclusion{
		BlendMode: BlendMode{
			name: "exclusion",
		},
	}
}

// Blend implements the exclusion blend mode
func (e *Exclusion) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform exclusion blend: a + b - 2ab
	r := exclusionBlend(baseLinear.R, blendLinear.R)
	g := exclusionBlend(baseLinear.G, blendLinear.G)
	b := exclusionBlend(baseLinear.B, blendLinear.B)
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

// exclusionBlend performs the exclusion blend operation for a single channel
func exclusionBlend(base, blend float64) float64 {
	return base + blend - 2*base*blend
}

// Register the blend mode
func init() {
	Register("exclusion", "Creates an effect similar to Difference but with lower contrast", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform exclusion blend: a + b - 2ab
		r := exclusionBlend(bottomLinear.R, topLinear.R)
		g := exclusionBlend(bottomLinear.G, topLinear.G)
		b := exclusionBlend(bottomLinear.B, topLinear.B)
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
