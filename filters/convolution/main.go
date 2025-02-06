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
				c := src.GetRGBA(math.Clamp(x+j-halfSize, 0, w-1), math.Clamp(y+i-halfSize, 0, h-1))
				weight := cm.Matrix[i][j]
				r += (float64(c.R()) / 255.0) * weight
				g += (float64(c.G()) / 255.0) * weight
				b += (float64(c.B()) / 255.0) * weight
			}
		}
		return x, y, rgba.New(
			math.Clamp((r*cm.Factor+cm.Bias)*255.0, 0.0, 255.0),
			math.Clamp((g*cm.Factor+cm.Bias)*255.0, 0.0, 255.0),
			math.Clamp((b*cm.Factor+cm.Bias)*255.0, 0.0, 255.0),
			float64(col.A()),
		)
	})
}

func (cm *ConvolutionMatrix) Apply3x3(src *image.Image) *image.Image {
	// assume kernel is 3x3
	w, h := src.W(), src.H()

	// Create the destination image
	dst := image.NewWithColor(w, h, *rgba.New(0, 0, 0, 0xFF))

	// Process row by row.
	for y := 0; y < h; y++ {
		// Precompute the clamped y indices for the 3 kernel rows.
		y0 := y - 1
		if y0 < 0 {
			y0 = 0
		}
		y1 := y
		y2 := y + 1
		if y2 >= h {
			y2 = h - 1
		}
		for x := 0; x < w; x++ {
			// Precompute the clamped x indices for the 3 kernel columns.
			x0 := x - 1
			if x0 < 0 {
				x0 = 0
			}
			x1 := x
			x2 := x + 1
			if x2 >= w {
				x2 = w - 1
			}

			// Get the neighboring pixels.
			c00 := src.GetRGBA(x0, y0)
			c01 := src.GetRGBA(x1, y0)
			c02 := src.GetRGBA(x2, y0)

			c10 := src.GetRGBA(x0, y1)
			c11 := src.GetRGBA(x1, y1)
			c12 := src.GetRGBA(x2, y1)

			c20 := src.GetRGBA(x0, y2)
			c21 := src.GetRGBA(x1, y2)
			c22 := src.GetRGBA(x2, y2)

			// Instead of a loop we explicitly calculate the convolution
			// Note: We assume that GetRGBA returns a value that you can query
			// for R(), G(), B() as integers 0-255.
			var r, g, b float64

			// Top row
			r += (float64(c00.R()) / 255.0) * cm.Matrix[0][0]
			g += (float64(c00.G()) / 255.0) * cm.Matrix[0][0]
			b += (float64(c00.B()) / 255.0) * cm.Matrix[0][0]

			r += (float64(c01.R()) / 255.0) * cm.Matrix[0][1]
			g += (float64(c01.G()) / 255.0) * cm.Matrix[0][1]
			b += (float64(c01.B()) / 255.0) * cm.Matrix[0][1]

			r += (float64(c02.R()) / 255.0) * cm.Matrix[0][2]
			g += (float64(c02.G()) / 255.0) * cm.Matrix[0][2]
			b += (float64(c02.B()) / 255.0) * cm.Matrix[0][2]

			// Middle row
			r += (float64(c10.R()) / 255.0) * cm.Matrix[1][0]
			g += (float64(c10.G()) / 255.0) * cm.Matrix[1][0]
			b += (float64(c10.B()) / 255.0) * cm.Matrix[1][0]

			r += (float64(c11.R()) / 255.0) * cm.Matrix[1][1]
			g += (float64(c11.G()) / 255.0) * cm.Matrix[1][1]
			b += (float64(c11.B()) / 255.0) * cm.Matrix[1][1]

			r += (float64(c12.R()) / 255.0) * cm.Matrix[1][2]
			g += (float64(c12.G()) / 255.0) * cm.Matrix[1][2]
			b += (float64(c12.B()) / 255.0) * cm.Matrix[1][2]

			// Bottom row
			r += (float64(c20.R()) / 255.0) * cm.Matrix[2][0]
			g += (float64(c20.G()) / 255.0) * cm.Matrix[2][0]
			b += (float64(c20.B()) / 255.0) * cm.Matrix[2][0]

			r += (float64(c21.R()) / 255.0) * cm.Matrix[2][1]
			g += (float64(c21.G()) / 255.0) * cm.Matrix[2][1]
			b += (float64(c21.B()) / 255.0) * cm.Matrix[2][1]

			r += (float64(c22.R()) / 255.0) * cm.Matrix[2][2]
			g += (float64(c22.G()) / 255.0) * cm.Matrix[2][2]
			b += (float64(c22.B()) / 255.0) * cm.Matrix[2][2]

			// Apply factor and bias
			R := math.Clamp((r*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)
			G := math.Clamp((g*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)
			B := math.Clamp((b*cm.Factor+cm.Bias)*255.0, 0.0, 255.0)

			// Set the output pixel. We assume the alpha channel remains unchanged.
			dst.SetRGBA(x, y, rgba.New(R, G, B, float64(c11.A())))
		}
	}
	src.Set(dst.Get())
	return src
}
