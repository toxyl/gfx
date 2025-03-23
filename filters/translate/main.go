// filters/translate/main.go
package translate

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("translate", []*meta.FilterMetaDataArg{
	{Name: "x", Default: 0.0},
	{Name: "y", Default: 0.0},
})

// Apply translates the image by the specified percentages (-1..1).
// It converts the percentages to absolute pixel offsets and applies the translation with wrap turned off.
func Apply(img *image.Image, x, y float64) *image.Image {
	absX := int(x * float64(img.W()))
	absY := int(y * float64(img.H()))
	img.Set(img.Translate(absX, absY, false).Get())
	return img
}
