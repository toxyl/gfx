package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// filters/transform/main.go

// TransformFunction represents a function that applies a 2D transformation to an image
type TransformFunction struct {
	*fx.BaseFunction
	matrix [3][3]float64 // 3x3 transformation matrix
}

// Function arguments
var transformArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "transform adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewTransform(matrix [3][3]float64) *TransformFunction {
	return &TransformFunction{
		BaseFunction: fx.NewBaseFunction("transform", "Applies transform transformation to an image", color.New(0, 0, 0, 1), transformArgs),
		matrix:       matrix,
	}
}

// Apply implements the Function interface
func (f *TransformFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Calculate new bounds after transformation
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)

	// Transform corners to find new bounds
	corners := [][2]float64{
		{float64(bounds.Min.X), float64(bounds.Min.Y)},
		{float64(bounds.Max.X), float64(bounds.Min.Y)},
		{float64(bounds.Min.X), float64(bounds.Max.Y)},
		{float64(bounds.Max.X), float64(bounds.Max.Y)},
	}

	for _, corner := range corners {
		x, y := f.transformPoint(corner[0], corner[1])
		minX = math.Min(minX, x)
		minY = math.Min(minY, y)
		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}

	// Create new image with transformed bounds
	newBounds := image.Rect(
		int(math.Floor(minX)),
		int(math.Floor(minY)),
		int(math.Ceil(maxX)),
		int(math.Ceil(maxY)),
	)
	dst := image.NewRGBA(newBounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Process each pixel in the destination image
	for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
		for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
			// Calculate inverse transformation
			srcX, srcY := f.inverseTransformPoint(float64(x), float64(y))

			// Skip if source point is outside original bounds
			if srcX < float64(bounds.Min.X) || srcX >= float64(bounds.Max.X) ||
				srcY < float64(bounds.Min.Y) || srcY >= float64(bounds.Max.Y) {
				continue
			}

			// Get the four surrounding pixels
			x0 := int(math.Floor(srcX))
			y0 := int(math.Floor(srcY))
			x1 := x0 + 1
			y1 := y0 + 1

			// Clamp coordinates
			if x1 >= bounds.Max.X {
				x1 = bounds.Max.X - 1
			}
			if y1 >= bounds.Max.Y {
				y1 = bounds.Max.Y - 1
			}

			// Calculate interpolation weights
			dx := srcX - float64(x0)
			dy := srcY - float64(y0)

			// Get colors of surrounding pixels
			c00 := temp.At(x0, y0)
			c10 := temp.At(x1, y0)
			c01 := temp.At(x0, y1)
			c11 := temp.At(x1, y1)

			// Interpolate
			r00, g00, b00, a00 := c00.RGBA()
			r10, g10, b10, a10 := c10.RGBA()
			r01, g01, b01, a01 := c01.RGBA()
			r11, g11, b11, a11 := c11.RGBA()

			// Perform bilinear interpolation
			r := uint8((float64(r00)*(1-dx)*(1-dy) + float64(r10)*dx*(1-dy) + float64(r01)*(1-dx)*dy + float64(r11)*dx*dy) / 65535.0)
			g := uint8((float64(g00)*(1-dx)*(1-dy) + float64(g10)*dx*(1-dy) + float64(g01)*(1-dx)*dy + float64(g11)*dx*dy) / 65535.0)
			b := uint8((float64(b00)*(1-dx)*(1-dy) + float64(b10)*dx*(1-dy) + float64(b01)*(1-dx)*dy + float64(b11)*dx*dy) / 65535.0)
			a := uint8((float64(a00)*(1-dx)*(1-dy) + float64(a10)*dx*(1-dy) + float64(a01)*(1-dx)*dy + float64(a11)*dx*dy) / 65535.0)

			dst.Set(x, y, stdcolor.RGBA{r, g, b, a})
		}
	}

	return dst, nil
}

// transformPoint applies the transformation matrix to a point
func (f *TransformFunction) transformPoint(x, y float64) (float64, float64) {
	// Apply transformation matrix
	newX := f.matrix[0][0]*x + f.matrix[0][1]*y + f.matrix[0][2]
	newY := f.matrix[1][0]*x + f.matrix[1][1]*y + f.matrix[1][2]
	w := f.matrix[2][0]*x + f.matrix[2][1]*y + f.matrix[2][2]

	// Apply perspective division
	if w != 0 {
		newX /= w
		newY /= w
	}

	return newX, newY
}

// inverseTransformPoint applies the inverse transformation matrix to a point
func (f *TransformFunction) inverseTransformPoint(x, y float64) (float64, float64) {
	// Calculate determinant
	det := f.matrix[0][0]*(f.matrix[1][1]*f.matrix[2][2]-f.matrix[1][2]*f.matrix[2][1]) -
		f.matrix[0][1]*(f.matrix[1][0]*f.matrix[2][2]-f.matrix[1][2]*f.matrix[2][0]) +
		f.matrix[0][2]*(f.matrix[1][0]*f.matrix[2][1]-f.matrix[1][1]*f.matrix[2][0])

	if det == 0 {
		return x, y // Return original point if matrix is singular
	}

	// Calculate inverse matrix
	inv := [3][3]float64{
		{
			(f.matrix[1][1]*f.matrix[2][2] - f.matrix[1][2]*f.matrix[2][1]) / det,
			(f.matrix[0][2]*f.matrix[2][1] - f.matrix[0][1]*f.matrix[2][2]) / det,
			(f.matrix[0][1]*f.matrix[1][2] - f.matrix[0][2]*f.matrix[1][1]) / det,
		},
		{
			(f.matrix[1][2]*f.matrix[2][0] - f.matrix[1][0]*f.matrix[2][2]) / det,
			(f.matrix[0][0]*f.matrix[2][2] - f.matrix[0][2]*f.matrix[2][0]) / det,
			(f.matrix[0][2]*f.matrix[1][0] - f.matrix[0][0]*f.matrix[1][2]) / det,
		},
		{
			(f.matrix[1][0]*f.matrix[2][1] - f.matrix[1][1]*f.matrix[2][0]) / det,
			(f.matrix[0][1]*f.matrix[2][0] - f.matrix[0][0]*f.matrix[2][1]) / det,
			(f.matrix[0][0]*f.matrix[1][1] - f.matrix[0][1]*f.matrix[1][0]) / det,
		},
	}

	// Apply inverse transformation
	srcX := inv[0][0]*x + inv[0][1]*y + inv[0][2]
	srcY := inv[1][0]*x + inv[1][1]*y + inv[1][2]
	w := inv[2][0]*x + inv[2][1]*y + inv[2][2]

	// Apply perspective division
	if w != 0 {
		srcX /= w
		srcY /= w
	}

	return srcX, srcY
}

func init() {
	fx.DefaultRegistry.Register(NewTransform([3][3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}))
}
