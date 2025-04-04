package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Color blend mode preserves the hue and saturation of the base color while using the luminosity of the blend color.
// This is useful for colorizing grayscale images or changing the color of an image while preserving its details.
//
// Formula: (base.hue, base.saturation, blend.luminosity)
func Color(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
	// Convert colors to float64 in range [0, 1]
	baseRf := float64(bottom.R) / 65535.0
	baseGf := float64(bottom.G) / 65535.0
	baseBf := float64(bottom.B) / 65535.0
	blendRf := float64(top.R) / 65535.0
	blendGf := float64(top.G) / 65535.0
	blendBf := float64(top.B) / 65535.0

	// Convert to HSL
	baseH, baseS, _ := rgbToHSL(baseRf, baseGf, baseBf)
	_, _, blendL := rgbToHSL(blendRf, blendGf, blendBf)

	// Convert back to RGB using base hue and saturation with blend luminosity
	r, g, b := hslToRGB(baseH, baseS, blendL)

	// Calculate alpha
	a := float64(bottom.A) * alpha

	result := &color.RGBA64{}
	result.R = r
	result.G = g
	result.B = b
	result.A = a
	return result
}

// rgbToHSL converts RGB values to HSL
func rgbToHSL(r, g, b float64) (h, s, l float64) {
	max := max(r, max(g, b))
	min := min(r, min(g, b))
	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
		return
	}

	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}
	h /= 6
	return
}

// hslToRGB converts HSL values to RGB
func hslToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		r = l
		g = l
		b = l
		return
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRGB(p, q, h+1.0/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1.0/3.0)
	return
}

// hueToRGB is a helper function for hslToRGB
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func init() {
	Register(constants.ModeColor, "color", constants.CategoryComponent, Color)
}
