package pastelize

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("pastelize", []*meta.FilterMetaDataArg{})

func Apply(i *image.Image) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		return x, y, col.SetS(col.S() * 0.5).SetL(col.L() + (1-col.L())*0.2)
	})
}
