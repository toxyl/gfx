package dialogbox

import (
	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/blendmodes/registry"
	"github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/types/blendmode"
)

func Draw(img *image.Image, x, y, w, h int, title string, content *image.Image, mode blendmode.BlendMode) error {
	width, height := content.Size()
	wc, hc := width, height
	hh := config.DIALOG_HEADER_HEIGHT
	x += config.SPACING
	y += config.SPACING

	// body
	err := img.DrawRect(x, y, w, h, 1, config.COLOR_BORDER, config.COLOR_FILL, registry.Get("normal"))
	if err != nil {
		return err
	}

	// shadow
	black := config.COLOR_BLACK_50_PCT
	for i := 0; i < config.SPACING; i++ {
		black.SetA(float64(config.SPACING-i) / float64(config.SPACING))
		err := img.DrawRect(x-i, y-i, w+i*2, h+i*2, 1, black, config.COLOR_TRANSPARENT, registry.Get("normal"))
		if err != nil {
			return err
		}
	}

	// content
	err = img.Blend(content, 0, 0, wc, hc, x+2, y+(h-hc)-2, wc, hc, mode, 1)
	if err != nil {
		return err
	}

	// header fill
	err = img.DrawRect(x, y, w, hh, 2, config.COLOR_BORDER_HEADER, config.COLOR_FILL_HEADER, registry.Get("normal"))
	if err != nil {
		return err
	}

	// title
	err = img.DrawText(title, x+config.SPACING, y+config.SPACING, *config.COLOR_FONT, false, registry.Get("normal"))
	if err != nil {
		return err
	}

	// content border
	err = img.DrawRect(x+1, y+hh, w-2, (h-hh)-1, 1, config.COLOR_BLACK_50_PCT, config.COLOR_TRANSPARENT, registry.Get("normal"))
	if err != nil {
		return err
	}

	return nil
}
