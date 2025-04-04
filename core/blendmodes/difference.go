package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Difference blend mode implementation
type Difference struct {
	BlendMode
}

// NewDifference creates a new Difference blend mode
func NewDifference() *Difference {
	return &Difference{
		BlendMode: BlendMode{
			name: "difference",
		},
	}
}

// Blend implements the difference blend mode
func (d *Difference) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform difference blend: |a - b|
	r := abs(baseLinear.R - blendLinear.R)
	g := abs(baseLinear.G - blendLinear.G)
	b := abs(baseLinear.B - blendLinear.B)
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

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Register the blend mode
func init() {
	Register("difference", "Subtracts the darker of the two colors from the lighter color", constants.CategoryComparative, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform difference blend: |a - b|
		r := abs(bottomLinear.R - topLinear.R)
		g := abs(bottomLinear.G - topLinear.G)
		b := abs(bottomLinear.B - topLinear.B)
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
