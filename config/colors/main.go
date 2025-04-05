package config

import (
	"github.com/toxyl/gfx/core/color"
)

// mustHSL creates a new HSL color and panics if there is an error
func mustHSL(h, s, l, a float64) *color.HSL {
	c, err := color.NewHSL(h, s, l, a)
	if err != nil {
		panic(err)
	}
	return c
}

var (
	COLOR_BORDER         = mustHSL(0, 0, 0.2, 0.85)
	COLOR_BORDER_HEADER  = mustHSL(0, 0, 0.3, 0.85)
	COLOR_FILL           = mustHSL(0, 0, 0.15, 0.85)
	COLOR_FILL_HEADER    = mustHSL(0, 0, 0.06, 0.85)
	COLOR_BLACK_12_5_PCT = mustHSL(0, 0, 0, 0.125)
	COLOR_BLACK_25_PCT   = mustHSL(0, 0, 0, 0.25)
	COLOR_BLACK_50_PCT   = mustHSL(0, 0, 0, 0.50)
	COLOR_TRANSPARENT    = mustHSL(0, 0, 0, 0)
	COLOR_FONT           = mustHSL(0, 0, 0.90, 1.0)
)
