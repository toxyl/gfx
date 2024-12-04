package blend

import (
	"github.com/toxyl/gfx/color/convert"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/math"
)

// RGBA blends two RGBA colors to a new RGBA color.
func RGBA(c1, c2 *rgba.RGBA, mode BlendMode, alpha float64) *rgba.RGBA {
	srcAlpha := float64(c1.A()) / 255.0
	dstAlpha := float64(c2.A()) / 255.0
	resAlpha := srcAlpha + dstAlpha*(1-srcAlpha) // Porter-Duff Over operator

	if resAlpha == 0 {
		return rgba.New(0, 0, 0, 0) // Return fully transparent black
	}

	c2.SetA(uint8(float64(c2.A()) * alpha))
	var blended *rgba.RGBA
	if blendFunc, found := BlendModes[mode]; found {
		blended = blendFunc(c1, c2)
	} else {
		blended = normal(c1, c2)
	}
	return rgba.New(blended.R(), blended.G(), blended.B(), uint8(math.Clamp(resAlpha*255.0, 0.0, 255.0)))
}

// HSLA blends two HSLA colors to a new HSLA color.
func HSLA(c1, c2 *hsla.HSLA, mode BlendMode, alpha float64) *hsla.HSLA {
	return convert.RGBAToHSLA(RGBA(convert.HSLAToRGBA(c1), convert.HSLAToRGBA(c2), mode, alpha))
}

// BlendMode type represents blend mode identifiers.
type BlendMode string

// BlendFunc type represents a blend mode function that blends two RGBA colors.
type BlendFunc func(c1, c2 *rgba.RGBA) *rgba.RGBA

// Constants for blend modes
const (
	NORMAL BlendMode = "normal"
	//
	DARKEN      BlendMode = "darken"
	MULTIPLY    BlendMode = "multiply"
	COLOR_BURN  BlendMode = "color-burn"
	LINEAR_BURN BlendMode = "linear-burn"
	//
	LIGHTEN BlendMode = "lighten"
	SCREEN  BlendMode = "screen"
	ADD     BlendMode = "add"
	//
	OVERLAY    BlendMode = "overlay"
	SOFT_LIGHT BlendMode = "soft-light"
	HARD_LIGHT BlendMode = "hard-light"
	PIN_LIGHT  BlendMode = "pin-light"
	//
	DIFFERENCE BlendMode = "difference"
	EXCLUSION  BlendMode = "exclusion"
	SUBTRACT   BlendMode = "subtract"
	DIVIDE     BlendMode = "divide"
	//
	AVERAGE  BlendMode = "average"
	NEGATION BlendMode = "negation"
	//
	ERASE BlendMode = "erase"
)

// BlendModes map allows accessing blend modes by name.
var BlendModes = map[BlendMode]BlendFunc{
	NORMAL: normal,
	//
	DARKEN:      darken,
	MULTIPLY:    multiply,
	COLOR_BURN:  colorBurn,
	LINEAR_BURN: linearBurn,
	//
	LIGHTEN: lighten,
	SCREEN:  screen,
	ADD:     add,
	//
	OVERLAY:    overlay,
	SOFT_LIGHT: softLight,
	HARD_LIGHT: hardLight,
	PIN_LIGHT:  pinLight,
	//
	DIFFERENCE: difference,
	EXCLUSION:  exclusion,
	SUBTRACT:   subtract,
	DIVIDE:     divide,
	//
	AVERAGE:  average,
	NEGATION: negation,
	//
	ERASE: erase,
}

// blendChannel is a reusable function to blend each color channel.
func blendChannel(src, dst uint8, a2 float64) uint8 {
	vSrc := float64(src) * (1 - a2)
	vDst := float64(dst) * a2
	return uint8(math.Round(math.Clamp(vSrc+vDst, 0x00, 0xFF)))
}

// blendModeFunc is a reusable function for applying a blend mode on each channel.
func blendModeFunc(c1, c2 *rgba.RGBA, blendFunc func(uint8, uint8) float64) *rgba.RGBA {
	a2 := float64(c2.A()) / 255.0

	if a2 == 0 {
		return c1
	}

	outR := uint8(math.Clamp(math.Round(float64(c1.R())*(1-a2)+blendFunc(c1.R(), c2.R())*a2), 0x00, 0xFF))
	outG := uint8(math.Clamp(math.Round(float64(c1.G())*(1-a2)+blendFunc(c1.G(), c2.G())*a2), 0x00, 0xFF))
	outB := uint8(math.Clamp(math.Round(float64(c1.B())*(1-a2)+blendFunc(c1.B(), c2.B())*a2), 0x00, 0xFF))

	return rgba.New(outR, outG, outB, 0xFF)
}

// normal implements the "normal" blend mode with cumulative alpha.
func normal(c1, c2 *rgba.RGBA) *rgba.RGBA {
	a2 := float64(c2.A()) / 255.0

	if a2 == 0 {
		return c1
	}

	outR := blendChannel(c1.R(), c2.R(), a2)
	outG := blendChannel(c1.G(), c2.G(), a2)
	outB := blendChannel(c1.B(), c2.B(), a2)

	return rgba.New(outR, outG, outB, 0xFF)
}

// multiply implements the "multiply" blend mode.
func multiply(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return (vSrc * vDst) / 255.0
	})
}

// lighten implements the "lighten" blend mode, where the lighter value of each color channel is selected.
func lighten(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return math.Max(vSrc, vDst)
	})
}

// darken implements the "darken" blend mode, where the darker value of each color channel is selected.
func darken(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return math.Min(vSrc, vDst)
	})
}

// screen implements the "screen" blend mode.
func screen(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return 255.0 * (1.0 - (1.0-vSrc/255.0)*(1.0-vDst/255.0))
	})
}

// add implements the "add" blend mode.
func add(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return math.Min(vSrc+vDst, 255.0)
	})
}

// overlay implements the "overlay" blend mode with smoother transitions.
func overlay(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src) / 255.0 // Normalize to [0, 1]
		vDst := float64(dst) / 255.0

		// Gamma-adjusted overlay blending for smoother transitions
		if vSrc < 0.5 {
			// Adjusted multiply for darker base colors
			return 2 * vSrc * vDst * 255.0
		}
		// Adjusted screen for lighter base colors
		return 255.0 * (1 - 2*(1-vSrc)*(1-vDst))
	})
}

// exclusion implements the "exclusion" blend mode.
func exclusion(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return (vSrc + vDst) - 2.0*vSrc*vDst/255.0
	})
}

// negation implements the "negation" blend mode.
func negation(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return (255.0 - math.Abs(255.0-vSrc-vDst))
	})
}

// colorBurn implements the "color burn" blend mode.
func colorBurn(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src) / 255.0
		vDst := float64(dst) / 255.0
		if vDst == 0 {
			return 0
		}
		return 255.0 * (1.0 - math.Min(1.0, (1.0-vSrc)/vDst))
	})
}

// linearBurn implements the "linear burn" blend mode.
func linearBurn(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		return math.Max(0, vSrc+vDst-255.0)
	})
}

// softLight implements the "soft light" blend mode.
func softLight(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src) / 255.0
		vDst := float64(dst) / 255.0
		if vDst < 0.5 {
			return 255.0 * (vSrc - (1.0-2.0*vDst)*vSrc*(1.0-vSrc))
		}
		return 255.0 * (vSrc + (2.0*vDst-1.0)*(math.Sqrt(vSrc)-vSrc))
	})
}

// hardLight implements the "hard light" blend mode.
func hardLight(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vDst := float64(dst) / 255.0
		if vDst < 0.5 {
			return 2.0 * float64(src) * vDst
		}
		return 255.0 - 2.0*(255.0-float64(src))*(1.0-vDst)
	})
}

// pinLight implements the "pin light" blend mode.
func pinLight(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		if vDst < 128 {
			return math.Min(vSrc, 2.0*vDst)
		}
		return math.Max(vSrc, 2.0*(vDst-128.0))
	})
}

// difference implements the "difference" blend mode.
func difference(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		return math.Abs(float64(src) - float64(dst))
	})
}

// subtract implements the "subtract" blend mode.
func subtract(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		return math.Max(0, float64(src)-float64(dst))
	})
}

// divide implements the "divide" blend mode.
func divide(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		vSrc := float64(src)
		vDst := float64(dst)
		if vDst == 0 {
			return 255
		}
		return math.Min(255.0, vSrc*255.0/vDst)
	})
}

// average implements the "average" blend mode.
func average(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return blendModeFunc(c1, c2, func(src, dst uint8) float64 {
		return (float64(src) + float64(dst)) / 2.0
	})
}

// erase implements the "erase" blend mode.
func erase(c1, c2 *rgba.RGBA) *rgba.RGBA {
	return rgba.New(c1.R(), c1.G(), c1.B(), uint8(math.Max(0, float64(c1.A())-float64(c2.A()))))
}
