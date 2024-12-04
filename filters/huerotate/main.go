package huerotate

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func Apply(i *image.Image, angle float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		newHue := math.Mod(col.H()+angle, 360)
		if newHue < 0 {
			newHue += 360
		}
		return x, y, col.SetH(newHue)
	})
}
