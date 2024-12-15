package convolution

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

var Meta = meta.New("convolution", []*meta.FilterMetaDataArg{
	{Name: "amount", Default: 1.0},
	{Name: "bias", Default: 0.0},
	{Name: "factor", Default: 1.0},
	{Name: "matrix", Default: [][]float64{
		{1.0, 1.0, 1.0},
		{1.0, 8.0, 1.0},
		{1.0, 1.0, 1.0},
	}},
})

type FilterFn func(intensity float64) (matrix [][]float64)

type ConvolutionMatrix struct {
	Matrix [][]float64
	Factor float64
	Bias   float64
}

func NewConvolutionMatrix(matrix [][]float64, factor, bias float64) *ConvolutionMatrix {
	return &ConvolutionMatrix{
		Matrix: matrix,
		Factor: factor,
		Bias:   bias,
	}
}

func NewCustomFilter(amount, factor, bias float64, filterFn FilterFn) *ConvolutionMatrix {
	return NewConvolutionMatrix(filterFn(amount), factor, bias)
}

func (cm *ConvolutionMatrix) Apply(src *image.Image) *image.Image {
	matrixSize := len(cm.Matrix)
	halfSize := matrixSize / 2
	w := src.W()
	h := src.H()

	return src.ProcessRGBA(0, 0, w, h, func(x, y int, col *rgba.RGBA) (x2, y2 int, col2 *rgba.RGBA) {
		var r, g, b float64

		for i := 0; i < matrixSize; i++ {
			for j := 0; j < matrixSize; j++ {
				px := math.Clamp(x+j-halfSize, 0, w-1)
				py := math.Clamp(y+i-halfSize, 0, h-1)

				c := src.GetRGBA(px, py)

				prf := float64(c.R()) / 255.0
				pgf := float64(c.G()) / 255.0
				pbf := float64(c.B()) / 255.0

				weight := cm.Matrix[i][j]
				r += prf * weight
				g += pgf * weight
				b += pbf * weight
			}
		}

		r = math.Clamp((r*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)
		g = math.Clamp((g*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)
		b = math.Clamp((b*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)

		return x, y, rgba.New(r, g, b, float64(col.A()))
	})
}
