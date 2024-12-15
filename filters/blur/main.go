package blur

import (
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

var Meta = meta.New("blur", []*meta.FilterMetaDataArg{
	{Name: "amount", Default: 1.0},
})

func Apply(img *image.Image, amount float64) *image.Image {
	kernelSize := math.Abs(int(2*amount)) + 1
	if kernelSize%2 == 0 {
		kernelSize++
	}
	matrix := make([][]float64, kernelSize)
	for i := range matrix {
		matrix[i] = make([]float64, kernelSize)
		for j := range matrix[i] {
			matrix[i][j] = 1.0 / float64(kernelSize*kernelSize)
		}
	}
	return convolution.NewConvolutionMatrix(matrix, 1, 0).Apply(img)
}
