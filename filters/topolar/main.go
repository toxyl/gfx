package topolar

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("to-polar", []*meta.FilterMetaDataArg{
	{Name: "angle-start", Default: 0.0},
	{Name: "angle-end", Default: 360.0},
	{Name: "rotation", Default: 0.0},
})

func Apply(img *image.Image, angleStart, angleEnd, rotation float64) *image.Image {
	img.Set(img.ToPolar(angleStart, angleEnd, rotation).Get())
	return img
}
