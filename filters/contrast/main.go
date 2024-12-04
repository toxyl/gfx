package contrast

import (
	"math"

	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image, contrast float64) *image.Image {
	if contrast < -1 || contrast > 1 {
		panic("Contrast value must be between -1 and 1")
	}

	// Precompute scaling factor
	factor := (259 * (contrast + 1)) / (255 * (1 - contrast))

	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		adjust := func(value uint8) uint8 {
			return uint8(math.Min(255, math.Max(0, factor*(float64(value)-128)+128)))
		}

		return x, y, rgba.New(
			adjust(col.R()),
			adjust(col.G()),
			adjust(col.B()),
			col.A(),
		)
	})
}
