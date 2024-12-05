package vars

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
)

var (
	CHAR_W                 = 9
	CHAR_H                 = 10
	SPRITESHEET_COLS       = 16
	DIALOG_HEADER_HEIGHT   = 18
	SPACING                = 4
	COLOR_BORDER           = hsla.New(0, 0, 0.2, 0.85)
	COLOR_BORDER_HEADER    = hsla.New(0, 0, 0.3, 0.85)
	COLOR_FILL             = hsla.New(0, 0, 0.15, 0.85)
	COLOR_FILL_HEADER      = hsla.New(0, 0, 0.06, 0.85)
	COLOR_BLACK_12_5_PCT   = hsla.New(0, 0, 0, 0.125)
	COLOR_BLACK_25_PCT     = hsla.New(0, 0, 0, 0.25)
	COLOR_BLACK_50_PCT     = hsla.New(0, 0, 0, 0.50)
	COLOR_TRANSPARENT      = hsla.New(0, 0, 0, 0)
	COLOR_TRANSPARENT_RGBA = rgba.New(0, 0, 0, 0)
	COLOR_FONT             = hsla.New(0, 0, 0.90, 1.0)
)
