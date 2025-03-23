package crop

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("crop", []*meta.FilterMetaDataArg{
	{Name: "left", Default: 0.0},
	{Name: "right", Default: 0.0},
	{Name: "top", Default: 0.0},
	{Name: "bottom", Default: 0.0},
})

func Apply(img *image.Image, left, right, top, bottom float64) *image.Image {
	w := float64(img.W())
	h := float64(img.H())
	x := left * w
	y := top * h
	x2 := (1 - right) * w
	y2 := (1 - bottom) * h
	img.Set(img.Crop(int(x), int(y), int(x2-x), int(y2-y), false).Get())
	return img
}
