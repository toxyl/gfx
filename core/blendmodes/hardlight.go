package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// HardLight blend mode implementation
type HardLight struct {
	BlendMode
}

// NewHardLight creates a new HardLight blend mode
func NewHardLight() *HardLight {
	return &HardLight{
		BlendMode: BlendMode{
			name: "hardlight",
		},
	}
}

// Blend implements the hard light blend mode
func (hl *HardLight) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform hard light blend
	r := hardLightBlend(baseLinear.R, blendLinear.R)
	g := hardLightBlend(baseLinear.G, blendLinear.G)
	b := hardLightBlend(baseLinear.B, blendLinear.B)
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

// hardLightBlend performs the hard light blend operation for a single channel
func hardLightBlend(base, blend float64) float64 {
	if blend <= 0.5 {
		return 2 * base * blend
	}
	return 1 - 2*(1-base)*(1-blend)
}

// Register the blend mode
func init() {
	Register("hardlight", "Applies a harsh lighting effect, similar to shining a harsh spotlight on the image", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform hard light blend
		r := hardLightBlend(bottomLinear.R, topLinear.R)
		g := hardLightBlend(bottomLinear.G, topLinear.G)
		b := hardLightBlend(bottomLinear.B, topLinear.B)
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
