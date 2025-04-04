package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// PinLight blend mode implementation
type PinLight struct {
	BlendMode
}

// NewPinLight creates a new PinLight blend mode
func NewPinLight() *PinLight {
	return &PinLight{
		BlendMode: BlendMode{
			name: "pinlight",
		},
	}
}

// Blend implements the pin light blend mode
func (pl *PinLight) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform pin light blend
	var (
		r     = pinLightBlend(baseLinear.R, blendLinear.R)
		g     = pinLightBlend(baseLinear.G, blendLinear.G)
		b     = pinLightBlend(baseLinear.B, blendLinear.B)
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

// pinLightBlend performs the pin light blend operation for a single channel
func pinLightBlend(base, blend float64) float64 {
	if blend <= 0.5 {
		// Darken: min(base, 2 * blend)
		return math.Min(base, 2*blend)
	}
	// Lighten: max(base, 2 * blend - 1)
	return math.Max(base, 2*blend-1)
}

// Register the blend mode
func init() {
	Register("pinlight", "Combines Darken and Lighten based on the blend color", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform pin light blend
		var (
			r = pinLightBlend(bottomLinear.R, topLinear.R)
			g = pinLightBlend(bottomLinear.G, topLinear.G)
			b = pinLightBlend(bottomLinear.B, topLinear.B)
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
