package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// RotateFunction represents a function that rotates an image
type RotateFunction struct {
	*fx.BaseFunction
	angle float64 // Rotation angle in degrees
}

// Function arguments
var rotateArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "rotate adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewRotate(angle float64) *RotateFunction {
	return &RotateFunction{
		BaseFunction: fx.NewBaseFunction("rotate", "Applies rotation transformation to an image", color.New(0, 0, 0, 1), rotateArgs),
		angle:        angle,
	}
}

// Apply implements the Function interface
func (f *RotateFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Convert angle to radians
	angle := f.angle * math.Pi / 180.0

	// Calculate new bounds after rotation
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
		x, y := f.rotatePoint(corner[0], corner[1], angle)
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
			// Calculate inverse rotation
			srcX, srcY := f.inverseRotatePoint(float64(x), float64(y), angle)

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

// rotatePoint applies rotation to a point
func (f *RotateFunction) rotatePoint(x, y, angle float64) (float64, float64) {
	sin, cos := math.Sin(angle), math.Cos(angle)
	return x*cos - y*sin, x*sin + y*cos
}

// inverseRotatePoint applies inverse rotation to a point
func (f *RotateFunction) inverseRotatePoint(x, y, angle float64) (float64, float64) {
	sin, cos := math.Sin(-angle), math.Cos(-angle)
	return x*cos - y*sin, x*sin + y*cos
}

func init() {
	fx.DefaultRegistry.Register(NewRotate(0))
}
