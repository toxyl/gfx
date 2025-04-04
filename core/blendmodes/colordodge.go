package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// ColorDodge blend mode implementation
type ColorDodge struct {
	BlendMode
}

// NewColorDodge creates a new ColorDodge blend mode
func NewColorDodge() *ColorDodge {
	return &ColorDodge{
		BlendMode: BlendMode{
			name: "colordodge",
		},
	}
}

// Blend implements the color dodge blend mode
func (cd *ColorDodge) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform color dodge blend: a / (1 - b)
	r := colorDodgeBlend(baseLinear.R, blendLinear.R)
	g := colorDodgeBlend(baseLinear.G, blendLinear.G)
	b := colorDodgeBlend(baseLinear.B, blendLinear.B)
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

// colorDodgeBlend performs the color dodge blend operation for a single channel
func colorDodgeBlend(base, blend float64) float64 {
	if blend == 1 {
		return 1
	}
	result := base / (1 - blend)
	if result > 1 {
		return 1
	}
	return result
}

// Register the blend mode
func init() {
	Register("colordodge", "Brightens the base color to reflect the blend color by decreasing contrast", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform color dodge blend: a / (1 - b)
		r := colorDodgeBlend(bottomLinear.R, topLinear.R)
		g := colorDodgeBlend(bottomLinear.G, topLinear.G)
		b := colorDodgeBlend(bottomLinear.B, topLinear.B)
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
