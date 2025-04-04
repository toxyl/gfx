package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// DarkerColor blend mode implementation
type DarkerColor struct {
	BlendMode
}

// NewDarkerColor creates a new DarkerColor blend mode
func NewDarkerColor() *DarkerColor {
	return &DarkerColor{
		BlendMode: BlendMode{
			name: "darkercolor",
		},
	}
}

// Blend implements the darker-color blend mode
func (dc *DarkerColor) Blend(base, blend *color.RGBA64) *color.RGBA64 {
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

	// Choose the color with lower luminance
	var (
		r     float64
		g     float64
		b     float64
		alpha = baseLinear.A
	)

	if baseLuminance <= blendLuminance {
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
	Register("darkercolor", "Selects the color with lower luminance between base and blend colors", constants.CategoryDarken, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
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

		// Choose the color with lower luminance
		var (
			r float64
			g float64
			b float64
			a = bottomLinear.A
		)

		if bottomLuminance <= topLuminance {
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
