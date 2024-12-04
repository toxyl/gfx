package vibrance

import (
	"math"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
)

func Apply(i *image.Image, boost float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		s := col.S()
		factor := 1 + boost*(1-s)
		newS := math.Min(1, s*factor)
		return x, y, col.SetS(newS)
	})
}
