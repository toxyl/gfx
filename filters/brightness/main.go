package brightness

import (
	"math"

	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image, adjustment float64) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		adjust := func(value uint8) uint8 {
			return uint8(math.Min(255, math.Max(0, float64(value)+adjustment)))
		}

		return x, y, rgba.New(
			adjust(col.R()),
			adjust(col.G()),
			adjust(col.B()),
			col.A(),
		)
	})
}
