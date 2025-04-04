package blendmodes

import (
	"math"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// VividLight blend mode implementation
type VividLight struct {
	BlendMode
}

// NewVividLight creates a new VividLight blend mode
func NewVividLight() *VividLight {
	return &VividLight{
		BlendMode: BlendMode{
			name: "vividlight",
		},
	}
}

// Blend implements the vivid light blend mode
func (vl *VividLight) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {
		// No additional processing needed
	})

	// Perform vivid light blend: b < 0.5 ? colorBurn : colorDodge
	var (
		r     = vividLightBlend(baseLinear.R, blendLinear.R)
		g     = vividLightBlend(baseLinear.G, blendLinear.G)
		b     = vividLightBlend(baseLinear.B, blendLinear.B)
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

// vividLightBlend performs the vivid light blend operation for a single channel
func vividLightBlend(base, blend float64) float64 {
	if blend <= 0.5 {
		// Color Burn: 1 - (1 - a) / (2 * b)
		if blend == 0 {
			return 0
		}
		result := 1 - (1-base)/(2*blend)
		return math.Max(0.0, result)
	}
	// Color Dodge: a / (2 * (1 - b))
	if blend == 1 {
		return 1
	}
	result := base / (2 * (1 - blend))
	return math.Min(1.0, result)
}

// Register the blend mode
func init() {
	Register("vividlight", "Combines Color Burn and Color Dodge based on the blend color", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {
			// No additional processing needed
		})

		// Perform vivid light blend
		var (
			r = vividLightBlend(bottomLinear.R, topLinear.R)
			g = vividLightBlend(bottomLinear.G, topLinear.G)
			b = vividLightBlend(bottomLinear.B, topLinear.B)
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
