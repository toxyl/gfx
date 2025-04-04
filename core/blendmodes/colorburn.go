package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// ColorBurn blend mode implementation
type ColorBurn struct {
	BlendMode
}

// NewColorBurn creates a new ColorBurn blend mode
func NewColorBurn() *ColorBurn {
	return &ColorBurn{
		BlendMode: BlendMode{
			name: "colorburn",
		},
	}
}

// Blend implements the color burn blend mode
func (cb *ColorBurn) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform color burn blend: 1 - (1 - a) / b
	r := colorBurnBlend(baseLinear.R, blendLinear.R)
	g := colorBurnBlend(baseLinear.G, blendLinear.G)
	b := colorBurnBlend(baseLinear.B, blendLinear.B)
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

// colorBurnBlend performs the color burn blend operation for a single channel
func colorBurnBlend(base, blend float64) float64 {
	if blend == 0 {
		return 0
	}
	result := 1 - (1-base)/blend
	if result < 0 {
		return 0
	}
	return result
}

// Register the blend mode
func init() {
	Register("colorburn", "Darkens the base color to reflect the blend color by increasing contrast", constants.CategoryDarken, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform color burn blend: 1 - (1 - a) / b
		r := colorBurnBlend(bottomLinear.R, topLinear.R)
		g := colorBurnBlend(bottomLinear.G, topLinear.G)
		b := colorBurnBlend(bottomLinear.B, topLinear.B)
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
