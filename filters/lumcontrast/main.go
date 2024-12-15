package lumcontrast

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

var Meta = meta.New("lum-contrast", []*meta.FilterMetaDataArg{
	{Name: "adjustment", Default: 0.0},
})

func Apply(i *image.Image, adjustment float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		l := col.L()
		newL := math.Clamp(0.5+(l-0.5)*(1.0-(-adjustment)), 0, 1)
		return x, y, col.SetL(newL)
	})
}
