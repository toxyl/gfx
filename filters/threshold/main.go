package threshold

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("threshold", []*meta.FilterMetaDataArg{
	{Name: "amount", Default: 0.0},
})

func Apply(img *image.Image, amount float64) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		gray := uint8(0.299*float64(col.R()) + 0.587*float64(col.G()) + 0.114*float64(col.B()))
		if gray > uint8(amount*255.0) {
			return x, y, col
		}
		return x, y, col.SetA(0)
	})
}
