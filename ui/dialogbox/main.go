package dialogbox

import (
	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/vars"
)

func Draw(img *image.Image, x, y, w, h int, title string, content *image.Image, mode blend.BlendMode) {
	wc := content.W()
	hc := content.H()
	hh := vars.DIALOG_HEADER_HEIGHT
	x += vars.SPACING
	y += vars.SPACING

	// body
	img.DrawRect(x, y, w, h, 1, vars.COLOR_BORDER, vars.COLOR_FILL, blend.NORMAL)

	// shadow
	black := vars.COLOR_BLACK_50_PCT
	for i := 0; i < vars.SPACING; i++ {
		black.SetA(float64(vars.SPACING-i) / float64(vars.SPACING))
		img.DrawRect(x-i, y-i, w+i*2, h+i*2, 1, black, vars.COLOR_TRANSPARENT, blend.NORMAL)
	}

	// content
	img.Draw(content, 0, 0, wc, hc, x+2, y+(h-hc)-2, wc, hc, mode, 1)

	// header fill
	img.DrawRect(x, y, w, hh, 2, vars.COLOR_BORDER_HEADER, vars.COLOR_FILL_HEADER, blend.NORMAL)

	// title
	img.DrawText(title, x+vars.SPACING, y+vars.SPACING, *vars.COLOR_FONT, false, blend.NORMAL)

	// content border
	img.DrawRect(x+1, y+hh, w-2, (h-hh)-1, 1, vars.COLOR_BLACK_50_PCT, vars.COLOR_TRANSPARENT, blend.NORMAL)
}
