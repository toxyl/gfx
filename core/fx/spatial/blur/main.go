package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// BlurFunction represents a function that applies a Gaussian blur to an image
type BlurFunction struct {
	*fx.BaseFunction
	radius int
}

// Function arguments
var blurArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "blur adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewBlur(radius int) *BlurFunction {
	return &BlurFunction{
		BaseFunction: fx.NewBaseFunction("blur", "Applies blur transformation to an image", color.New(0, 0, 0, 1), blurArgs),
		radius:       radius,
	}
}

// Apply implements the Function interface
func (f *BlurFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	// Create temporary image for horizontal blur
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Apply horizontal blur
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.applyHorizontalBlur(x, y, img)
			temp.Set(x, y, col.ToUint8())
		}
	}

	// Convert temp to image.Image for vertical blur
	var tempImg image.Image = temp

	// Apply vertical blur
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.applyVerticalBlur(x, y, &tempImg)
			dst.Set(x, y, col.ToUint8())
		}
	}

	result := image.Image(dst)
	return &result, nil
}

// applyHorizontalBlur applies Gaussian blur in the horizontal direction
func (f *BlurFunction) applyHorizontalBlur(x, y int, img *image.Image) *color.Color64 {
	var r, g, b, a float64
	var weightSum float64

	bounds := (*img).Bounds()
	kernel := f.createGaussianKernel()

	for i := -f.radius; i <= f.radius; i++ {
		px := x + i
		if px < bounds.Min.X || px >= bounds.Max.X {
			continue
		}

		col := color.New(0, 0, 0, 1)
		col.SetUint8((*img).At(px, y).(stdcolor.RGBA))

		weight := kernel[i+f.radius]
		r += col.R * weight
		g += col.G * weight
		b += col.B * weight
		a += col.A * weight
		weightSum += weight
	}

	return color.New(r/weightSum, g/weightSum, b/weightSum, a/weightSum)
}

// applyVerticalBlur applies Gaussian blur in the vertical direction
func (f *BlurFunction) applyVerticalBlur(x, y int, img *image.Image) *color.Color64 {
	var r, g, b, a float64
	var weightSum float64

	bounds := (*img).Bounds()
	kernel := f.createGaussianKernel()

	for i := -f.radius; i <= f.radius; i++ {
		py := y + i
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}

		col := color.New(0, 0, 0, 1)
		col.SetUint8((*img).At(x, py).(stdcolor.RGBA))

		weight := kernel[i+f.radius]
		r += col.R * weight
		g += col.G * weight
		b += col.B * weight
		a += col.A * weight
		weightSum += weight
	}

	return color.New(r/weightSum, g/weightSum, b/weightSum, a/weightSum)
}

// createGaussianKernel creates a 1D Gaussian kernel
func (f *BlurFunction) createGaussianKernel() []float64 {
	kernel := make([]float64, 2*f.radius+1)
	var sum float64

	sigma := float64(f.radius) / 2.0
	for i := -f.radius; i <= f.radius; i++ {
		x := float64(i)
		kernel[i+f.radius] = math.Exp(-(x*x)/(2*sigma*sigma)) / (math.Sqrt(2*math.Pi) * sigma)
		sum += kernel[i+f.radius]
	}

	// Normalize kernel
	for i := range kernel {
		kernel[i] /= sum
	}

	return kernel
}
