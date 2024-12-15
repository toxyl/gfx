package sharpen

import (
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("sharpen", []*meta.FilterMetaDataArg{
	{Name: "amount", Default: 0.0},
})

func Apply(img *image.Image, amount float64) *image.Image {
	return convolution.NewConvolutionMatrix([][]float64{
		{-amount, -amount, -amount},
		{-amount, 1 + 8*amount, -amount},
		{-amount, -amount, -amount},
	}, 1, 0).Apply(img)
}
