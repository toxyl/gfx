package fliph

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("flip-h", []*meta.FilterMetaDataArg{})

func Apply(img *image.Image) *image.Image {
	img.Set(img.FlipHorizontal().Get())
	return img
}
