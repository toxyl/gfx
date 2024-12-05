package sat

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func Apply(img *image.Image, saturation float64) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		r := float64(col.R())
		g := float64(col.G())
		b := float64(col.B())
		gray := 0.299*r + 0.587*g + 0.114*b
		adjust := func(value float64) uint8 {
			return uint8(math.Clamp(gray+(1-(-saturation))*(value-gray), 0.0, 255.0))
		}
		return x, y, rgba.New(adjust(r), adjust(g), adjust(b), col.A())
	})
}
