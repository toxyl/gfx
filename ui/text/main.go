package text

import (
	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/image"
)

func Draw(img *image.Image, x, y int, text string, mode blendmodes.IBlendMode) error {
	return img.DrawText(text, x, y, *config.COLOR_FONT, false, mode)
}
