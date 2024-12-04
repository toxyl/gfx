package text

import (
	"github.com/toxy/gfx/color/blend"
	"github.com/toxy/gfx/image"
	"github.com/toxy/gfx/vars"
)

func Draw(img *image.Image, x, y int, text string, mode blend.BlendMode) {
	img.DrawText(text, x, y, *vars.COLOR_FONT, false, mode)
}
