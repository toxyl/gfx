package hue

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
)

func Apply(i *image.Image, hue float64) *image.Image {
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) { return x, y, col.ShiftH(hue * 360) })
}