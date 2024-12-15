package vibrance

import (
	"math"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("vibrance", []*meta.FilterMetaDataArg{
	{Name: "adjustment", Default: 0.0},
})

func Apply(i *image.Image, adjustment float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		s := col.S()
		factor := 1 + adjustment*(1-s)
		newS := math.Min(1, s*factor)
		return x, y, col.SetS(newS)
	})
}
