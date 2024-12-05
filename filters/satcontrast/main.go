package satcontrast

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func Apply(i *image.Image, factor float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		s := col.S()
		newS := math.Clamp(0.5+(s-0.5)*(1.0-factor), 0, 1)
		return x, y, col.SetS(newS)
	})
}
