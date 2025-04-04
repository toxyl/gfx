package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// LinearLight blend mode implementation
type LinearLight struct {
	BlendMode
}

// NewLinearLight creates a new LinearLight blend mode
func NewLinearLight() *LinearLight {
	return &LinearLight{
		BlendMode: BlendMode{
			name: "linearlight",
		},
	}
}

// Blend implements the linear light blend mode
func (ll *LinearLight) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform linear light blend: b < 0.5 ? linearBurn : linearDodge
	var (
		r     = linearLightBlend(baseLinear.R, blendLinear.R)
		g     = linearLightBlend(baseLinear.G, blendLinear.G)
		b     = linearLightBlend(baseLinear.B, blendLinear.B)
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

// linearLightBlend performs the linear light blend operation for a single channel
func linearLightBlend(base, blend float64) float64 {
	if blend <= 0.5 {
		// Linear Burn: a + 2*b - 1
		return math.Max(0.0, base+2*blend-1.0)
	}
	// Linear Dodge: a + 2*(b-0.5)
	return math.Min(1.0, base+2*(blend-0.5))
}

// Register the blend mode
func init() {
	Register("linearlight", "Combines Linear Dodge and Linear Burn", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform linear light blend
		var (
			r = linearLightBlend(bottomLinear.R, topLinear.R)
			g = linearLightBlend(bottomLinear.G, topLinear.G)
			b = linearLightBlend(bottomLinear.B, topLinear.B)
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
