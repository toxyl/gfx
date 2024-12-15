package brightness

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

var Meta = meta.New("brightness", []*meta.FilterMetaDataArg{
	{Name: "adjustment", Default: 1.0},
})

func Apply(img *image.Image, adjustment float64) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		adjust := func(value uint8) uint8 {
			return uint8(math.Clamp(float64(value)+adjustment*255.0, 0.0, 255.0))
		}
		return x, y, rgba.New(
			adjust(col.R()),
			adjust(col.G()),
			adjust(col.B()),
			col.A(),
		)
	})
}
