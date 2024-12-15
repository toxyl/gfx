package colorshift

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("color-shift", []*meta.FilterMetaDataArg{
	{Name: "hue", Default: 0.0},
	{Name: "sat", Default: 0.0},
	{Name: "lum", Default: 0.0},
})

func Apply(i *image.Image, hue, sat, lum float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, hsla *hsla.HSLA) (x2, y2 int, col2 *hsla.HSLA) {
		return x, y, hsla.Shift(hue, sat, lum, 0)
	})
}
