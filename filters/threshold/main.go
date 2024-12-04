package threshold

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image, threshold uint8) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		gray := uint8(0.299*float64(col.R()) + 0.587*float64(col.G()) + 0.114*float64(col.B()))
		if gray > threshold {
			return x, y, rgba.New(255, 255, 255, col.A())
		}
		return x, y, rgba.New(0, 0, 0, col.A())
	})
}
