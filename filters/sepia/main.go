package sepia

import (
	"math"

	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		r := float64(col.R())
		g := float64(col.G())
		b := float64(col.B())

		newR := math.Min(255, 0.393*r+0.769*g+0.189*b)
		newG := math.Min(255, 0.349*r+0.686*g+0.168*b)
		newB := math.Min(255, 0.272*r+0.534*g+0.131*b)

		return x, y, rgba.New(uint8(newR), uint8(newG), uint8(newB), col.A())
	})
}
