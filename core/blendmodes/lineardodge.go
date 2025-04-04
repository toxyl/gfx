package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

func init() {
	Register("lineardodge", "Brightens the base color by increasing brightness", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {})

		// Calculate final alpha
		finalAlpha := alphaComposite(bottomLinear, topLinear, alpha)

		// Perform linear dodge blend
		r := bottomLinear.R + topLinear.R
		g := bottomLinear.G + topLinear.G
		b := bottomLinear.B + topLinear.B

		// Clamp values to [0, 1]
		r = math.Min(1, math.Max(0, r))
		g = math.Min(1, math.Max(0, g))
		b = math.Min(1, math.Max(0, b))

		// Return result
		return &color.RGBA64{
			R: r,
			G: g,
			B: b,
			A: finalAlpha,
		}
	})
}
