package gamma

import (
	"math"

	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("gamma", []*meta.FilterMetaDataArg{
	{Name: "adjustment", Default: 1.0},
})

func Apply(img *image.Image, adjustment float64) *image.Image {
	// Precompute gamma correction lookup table
	lut := make([]uint8, 256)
	invGamma := 1.0 / (adjustment + 1)
	for i := range 256 {
		lut[i] = uint8(math.Min(255, math.Pow(float64(i)/255.0, invGamma)*255.0))
	}

	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		return x, y, rgba.New(lut[col.R()], lut[col.G()], lut[col.B()], col.A())
	})
}
