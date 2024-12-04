package lightnesscontrast

import (
	"math"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
)

func Apply(i *image.Image, factor float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		l := col.L()
		newL := math.Min(1, math.Max(0, 0.5+(l-0.5)*factor))
		return x, y, col.SetL(newL)
	})
}
