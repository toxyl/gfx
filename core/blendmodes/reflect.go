package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Reflect blend mode implementation
type Reflect struct {
	BlendMode
}

// NewReflect creates a new Reflect blend mode
func NewReflect() *Reflect {
	return &Reflect{
		BlendMode: BlendMode{
			name: "reflect",
		},
	}
}

// Blend implements the reflect blend mode
func (rf *Reflect) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform reflect blend
	var (
		r     = reflectBlend(baseLinear.R, blendLinear.R)
		g     = reflectBlend(baseLinear.G, blendLinear.G)
		b     = reflectBlend(baseLinear.B, blendLinear.B)
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

// reflectBlend performs the reflect blend operation for a single channel
func reflectBlend(base, blend float64) float64 {
	if blend == 0 {
		return 0
	}
	if blend == 1 {
		return 1
	}
	// Reflect: if blend == 1 then 1, else base * base / (1 - blend)
	result := base * base / (1 - blend)
	return math.Min(1.0, result)
}

// Register the blend mode
func init() {
	Register("reflect", "Creates a reflective effect based on the base and blend colors", constants.CategorySpecial, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform reflect blend
		var (
			r = reflectBlend(bottomLinear.R, topLinear.R)
			g = reflectBlend(bottomLinear.G, topLinear.G)
			b = reflectBlend(bottomLinear.B, topLinear.B)
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
