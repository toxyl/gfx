package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// SoftLight blend mode implementation
type SoftLight struct {
	BlendMode
}

// NewSoftLight creates a new SoftLight blend mode
func NewSoftLight() *SoftLight {
	return &SoftLight{
		BlendMode: BlendMode{
			name: "softlight",
		},
	}
}

// Blend implements the soft light blend mode
func (sl *SoftLight) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform soft light blend
	r := softLightBlend(baseLinear.R, blendLinear.R)
	g := softLightBlend(baseLinear.G, blendLinear.G)
	b := softLightBlend(baseLinear.B, blendLinear.B)
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

// softLightBlend performs the soft light blend operation for a single channel
func softLightBlend(base, blend float64) float64 {
	if blend <= 0.5 {
		return base - (1-2*blend)*base*(1-base)
	}
	return base + (2*blend-1)*(sqrt(base)-base)
}

// sqrt returns the square root of a float64 value
func sqrt(x float64) float64 {
	// Using Newton's method for square root approximation
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// Register the blend mode
func init() {
	Register("softlight", "Similar to Overlay but produces softer results", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform soft light blend
		r := softLightBlend(bottomLinear.R, topLinear.R)
		g := softLightBlend(bottomLinear.G, topLinear.G)
		b := softLightBlend(bottomLinear.B, topLinear.B)
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
