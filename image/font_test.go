package image

import (
	_ "embed"
	"testing"

	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
)

func TestImage_DrawText(t *testing.T) {
	var (
		fontColor = hsla.New(30, 1.0, 0.5, 1.0)
	)
	type test struct {
		name string
		text string
		x    int
		y    int
		col  *hsla.HSLA
		glow bool
		mode blend.BlendMode
	}
	tests := []test{
		{"Normal", "Hello World", 1, 2, fontColor, false, blend.NORMAL},
		{"Add", "Hello World", 1, 2, fontColor, false, blend.ADD},
		{"Darken", "Hello World", 1, 2, fontColor, false, blend.DARKEN},
		{"Overlay", "Hello World", 1, 2, fontColor, false, blend.OVERLAY},
		{"Screen", "Hello World", 1, 2, fontColor, false, blend.SCREEN},
		{"Exclusion", "Hello World", 1, 2, fontColor, false, blend.EXCLUSION},
		{"Lighten", "Hello World", 1, 2, fontColor, false, blend.LIGHTEN},
		{"Multiply", "Hello World", 1, 2, fontColor, false, blend.MULTIPLY},
		{"Negation", "Hello World", 1, 2, fontColor, false, blend.NEGATION},
		{"Average", "Hello World", 1, 2, fontColor, false, blend.AVERAGE},
		{"Color Burn", "Hello World", 1, 2, fontColor, false, blend.COLOR_BURN},
		{"Difference", "Hello World", 1, 2, fontColor, false, blend.DIFFERENCE},
		{"Divide", "Hello World", 1, 2, fontColor, false, blend.DIVIDE},
		{"Hard Light", "Hello World", 1, 2, fontColor, false, blend.HARD_LIGHT},
		{"Linear Burn", "Hello World", 1, 2, fontColor, false, blend.LINEAR_BURN},
		{"Pin Light", "Hello World", 1, 2, fontColor, false, blend.PIN_LIGHT},
		{"Soft Light", "Hello World", 1, 2, fontColor, false, blend.SOFT_LIGHT},
		{"Subtract", "Hello World", 1, 2, fontColor, false, blend.SUBTRACT},

		{"Normal - Glow", "Hello World", 1, 2, fontColor, true, blend.NORMAL},
		{"Add - Glow", "Hello World", 1, 2, fontColor, true, blend.ADD},
		{"Darken - Glow", "Hello World", 1, 2, fontColor, true, blend.DARKEN},
		{"Overlay - Glow", "Hello World", 1, 2, fontColor, true, blend.OVERLAY},
		{"Screen - Glow", "Hello World", 1, 2, fontColor, true, blend.SCREEN},
		{"Exclusion - Glow", "Hello World", 1, 2, fontColor, true, blend.EXCLUSION},
		{"Lighten - Glow", "Hello World", 1, 2, fontColor, true, blend.LIGHTEN},
		{"Multiply - Glow", "Hello World", 1, 2, fontColor, true, blend.MULTIPLY},
		{"Negation - Glow", "Hello World", 1, 2, fontColor, true, blend.NEGATION},
		{"Average - Glow", "Hello World", 1, 2, fontColor, true, blend.AVERAGE},
		{"Color Burn - Glow", "Hello World", 1, 2, fontColor, true, blend.COLOR_BURN},
		{"Difference - Glow", "Hello World", 1, 2, fontColor, true, blend.DIFFERENCE},
		{"Divide - Glow", "Hello World", 1, 2, fontColor, true, blend.DIVIDE},
		{"Hard Light - Glow", "Hello World", 1, 2, fontColor, true, blend.HARD_LIGHT},
		{"Linear Burn - Glow", "Hello World", 1, 2, fontColor, true, blend.LINEAR_BURN},
		{"Pin Light - Glow", "Hello World", 1, 2, fontColor, true, blend.PIN_LIGHT},
		{"Soft Light - Glow", "Hello World", 1, 2, fontColor, true, blend.SOFT_LIGHT},
		{"Subtract - Glow", "Hello World", 1, 2, fontColor, true, blend.SUBTRACT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewWithColor(101, 13, *rgba.New(0x00, 0x33, 0x33, 0xFF))
			i.DrawText(tt.text, tt.x, tt.y, *tt.col, tt.glow, tt.mode).SaveAsPNG("../test_data/text/" + tt.name + ".png")
		})
	}
}
