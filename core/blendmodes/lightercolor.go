package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// LighterColor blend mode implementation
type LighterColor struct {
	BlendMode
}

// NewLighterColor creates a new LighterColor blend mode
func NewLighterColor() *LighterColor {
	return &LighterColor{
		BlendMode: BlendMode{
			name: "lightercolor",
		},
	}
}

// Blend implements the lighter-color blend mode
func (lc *LighterColor) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Calculate luminance for both colors
	baseLuminance := 0.2126*baseLinear.R + 0.7152*baseLinear.G + 0.0722*baseLinear.B
	blendLuminance := 0.2126*blendLinear.R + 0.7152*blendLinear.G + 0.0722*blendLinear.B

	// Choose the color with higher luminance
	var (
		r     float64
		g     float64
		b     float64
		alpha = baseLinear.A
	)

	if baseLuminance >= blendLuminance {
		r = baseLinear.R
		g = baseLinear.G
		b = baseLinear.B
	} else {
		r = blendLinear.R
		g = blendLinear.G
		b = blendLinear.B
	}

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
	Register("lightercolor", "Selects the lighter of the base and blend colors", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Calculate luminance for both colors
		bottomLuminance := 0.2126*bottomLinear.R + 0.7152*bottomLinear.G + 0.0722*bottomLinear.B
		topLuminance := 0.2126*topLinear.R + 0.7152*topLinear.G + 0.0722*topLinear.B

		// Choose the color with higher luminance
		var (
			r float64
			g float64
			b float64
			a = bottomLinear.A
		)

		if bottomLuminance >= topLuminance {
			r = bottomLinear.R
			g = bottomLinear.G
			b = bottomLinear.B
		} else {
			r = topLinear.R
			g = topLinear.G
			b = topLinear.B
		}

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
