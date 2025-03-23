package cropcircle

import (
	"math"

	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("crop-circle", []*meta.FilterMetaDataArg{
	{Name: "radius", Default: 0.0},
	{Name: "offset-x", Default: 0.0},
	{Name: "offset-y", Default: 0.0},
})

func Apply(img *image.Image, radius, offsetX, offsetY float64) *image.Image {
	w, h, hw, hh := float64(img.W()), float64(img.H()), float64(img.CW()), float64(img.CH())
	img.Set(img.CropCircle(int(hw+offsetX*hw), int(hh+offsetY*hh), int(radius*math.Max(w, h)), false).Get())
	return img
}
