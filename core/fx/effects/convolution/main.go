package fx

import (
	"errors"
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// ErrInvalidMatrix is returned when the convolution matrix is invalid
var ErrInvalidMatrix = errors.New("invalid convolution matrix")

// ConvolutionMatrix represents a convolution matrix for image processing
type ConvolutionMatrix struct {
	Matrix  [][]float64
	Divisor float64
	Offset  float64
}

// NewConvolutionMatrix creates a new convolution matrix
func NewConvolutionMatrix(matrix [][]float64, divisor, offset float64) *ConvolutionMatrix {
	return &ConvolutionMatrix{
		Matrix:  matrix,
		Divisor: divisor,
		Offset:  offset,
	}
}

// ConvolutionFunction represents a function that applies a convolution matrix to an image
type ConvolutionFunction struct {
	*fx.BaseFunction
	matrix *ConvolutionMatrix
}

// Function arguments
var convolutionArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "convolution adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewConvolution(matrix *ConvolutionMatrix) *ConvolutionFunction {
	return &ConvolutionFunction{
		BaseFunction: fx.NewBaseFunction("convolution", "Applies convolution transformation to an image", color.New(0, 0, 0, 1), convolutionArgs),
		matrix:       matrix,
	}
}

// Apply implements the Function interface
func (f *ConvolutionFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}
			dst.Set(x, y, col.ToUint8())
		}
	}

	return dst, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *ConvolutionFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	var r, g, b, a float64
	var weightSum float64

	matrix := f.matrix.Matrix
	matrixSize := len(matrix)
	halfSize := matrixSize / 2

	for dy := -halfSize; dy <= halfSize; dy++ {
		for dx := -halfSize; dx <= halfSize; dx++ {
			px := x + dx
			py := y + dy

			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}

			col := color.New(0, 0, 0, 1)
			col.SetUint8(img.At(px, py).(stdcolor.RGBA))

			weight := matrix[dy+halfSize][dx+halfSize]
			r += col.R * weight
			g += col.G * weight
			b += col.B * weight
			a += col.A * weight
			weightSum += weight
		}
	}

	// Apply divisor and offset
	if f.matrix.Divisor != 0 {
		r = r/f.matrix.Divisor + f.matrix.Offset
		g = g/f.matrix.Divisor + f.matrix.Offset
		b = b/f.matrix.Divisor + f.matrix.Offset
		a = a/f.matrix.Divisor + f.matrix.Offset
	}

	// Clamp values
	r = clamp(r)
	g = clamp(g)
	b = clamp(b)
	a = clamp(a)

	return color.New(r, g, b, a), nil
}

// clamp ensures a value is between 0 and 1
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Apply3x3 applies a 3x3 convolution matrix to an image
func Apply3x3(img image.Image, matrix [][]float64, divisor, offset float64) (image.Image, error) {
	if len(matrix) != 3 || len(matrix[0]) != 3 {
		return nil, ErrInvalidMatrix
	}

	conv := NewConvolution(NewConvolutionMatrix(matrix, divisor, offset))
	return conv.Apply(img)
}

func init() {
	fx.DefaultRegistry.Register(NewConvolution(NewConvolutionMatrix([][]float64{
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}, 1, 0)))
}
