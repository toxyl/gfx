package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Lighten blend mode implementation
type Lighten struct {
	BlendMode
}

// NewLighten creates a new Lighten blend mode
func NewLighten() *Lighten {
	return &Lighten{
		BlendMode: BlendMode{
			name: "lighten",
		},
	}
}

// Blend implements the lighten blend mode
func (l *Lighten) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform lighten blend: max(a, b)
	r := max(baseLinear.R, blendLinear.R)
	g := max(baseLinear.G, blendLinear.G)
	b := max(baseLinear.B, blendLinear.B)
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

// max returns the maximum of two float64 values
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Register the blend mode
func init() {
	Register("lighten", "Selects the lighter of the base and blend colors for each channel", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform lighten blend: max(a, b)
		r := max(bottomLinear.R, topLinear.R)
		g := max(bottomLinear.G, topLinear.G)
		b := max(bottomLinear.B, topLinear.B)
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
