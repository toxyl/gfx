package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// HardMix blend mode implementation
type HardMix struct {
	BlendMode
}

// NewHardMix creates a new HardMix blend mode
func NewHardMix() *HardMix {
	return &HardMix{
		BlendMode: BlendMode{
			name: "hardmix",
		},
	}
}

// Blend implements the hard mix blend mode
func (hm *HardMix) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform hard mix blend
	var (
		r     = hardMixBlend(baseLinear.R, blendLinear.R)
		g     = hardMixBlend(baseLinear.G, blendLinear.G)
		b     = hardMixBlend(baseLinear.B, blendLinear.B)
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

// hardMixBlend performs the hard mix blend operation for a single channel
func hardMixBlend(base, blend float64) float64 {
	// Hard Mix: if (base + blend) > 1 then 1, else 0
	if base+blend > 1.0 {
		return 1.0
	}
	return 0.0
}

// Register the blend mode
func init() {
	Register("hardmix", "Extreme version of Linear Light that clips values to pure black or white", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform hard mix blend
		var (
			r = hardMixBlend(bottomLinear.R, topLinear.R)
			g = hardMixBlend(bottomLinear.G, topLinear.G)
			b = hardMixBlend(bottomLinear.B, topLinear.B)
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
