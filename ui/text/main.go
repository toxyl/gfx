package text

import (
	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/vars"
)

func Draw(img *image.Image, x, y int, text string, mode blend.BlendMode) {
	img.DrawText(text, x, y, *vars.COLOR_FONT, false, mode)
}
