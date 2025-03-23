package flipv

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("flip-v", []*meta.FilterMetaDataArg{})

func Apply(img *image.Image) *image.Image {
	img.Set(img.FlipVertical().Get())
	return img
}
