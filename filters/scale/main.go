// filters/scale/main.go
package scale

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("scale", []*meta.FilterMetaDataArg{
	{Name: "scale", Default: 0.0},
	{Name: "offset-x", Default: 0.0},
	{Name: "offset-y", Default: 0.0},
})

// Apply scales the image by the specified factor around a center defined by offsets (-1..1) from the image center.
func Apply(img *image.Image, scaleFactor, offsetX, offsetY float64) *image.Image {
	cx := int(float64(img.CW()) + offsetX*float64(img.CW()))
	cy := int(float64(img.CH()) + offsetY*float64(img.CH()))
	img.Set(img.Scale(scaleFactor, cx, cy).Get())
	return img
}
