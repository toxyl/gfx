package invert

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		return x, y, rgba.New(
			255-col.R(),
			255-col.G(),
			255-col.B(),
			col.A(),
		)
	})
}
