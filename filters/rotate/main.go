package rotate

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("rotate", []*meta.FilterMetaDataArg{
	{Name: "angle", Default: 0.0},
	{Name: "offset-x", Default: 0.0},
	{Name: "offset-y", Default: 0.0},
})

func Apply(img *image.Image, angle, offsetX, offsetY float64) *image.Image {
	hw := float64(img.CW())
	hh := float64(img.CH())
	img.Set(img.Rotate(angle, hw+offsetX*hw, hh+offsetY*hh).Get())
	return img
}
