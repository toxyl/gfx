package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// filters/scale/main.go

// ScaleFunction represents a function that scales an image
type ScaleFunction struct {
	*fx.BaseFunction
	scaleX, scaleY float64 // Scale factors
}

// Function arguments
var scaleArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "scale adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewScale(scaleX, scaleY float64) *ScaleFunction {
	return &ScaleFunction{
		BaseFunction: fx.NewBaseFunction("scale", "Applies scaling transformation to an image", color.New(0, 0, 0, 1), scaleArgs),
		scaleX:       math.Max(0.0, scaleX),
		scaleY:       math.Max(0.0, scaleY),
	}
}

// Apply implements the Function interface
func (f *ScaleFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Calculate new dimensions
	newWidth := int(float64(bounds.Dx()) * f.scaleX)
	newHeight := int(float64(bounds.Dy()) * f.scaleY)

	// Create new image with scaled dimensions
	newBounds := image.Rect(0, 0, newWidth, newHeight)
	dst := image.NewRGBA(newBounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Process each pixel in the destination image
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate source coordinates
			srcX := float64(x) / f.scaleX
			srcY := float64(y) / f.scaleY

			// Skip if source point is outside original bounds
			if srcX < 0 || srcX >= float64(bounds.Dx()) ||
				srcY < 0 || srcY >= float64(bounds.Dy()) {
				continue
			}

			// Get the four surrounding pixels
			x0 := int(math.Floor(srcX))
			y0 := int(math.Floor(srcY))
			x1 := x0 + 1
			y1 := y0 + 1

			// Clamp coordinates
			if x1 >= bounds.Dx() {
				x1 = bounds.Dx() - 1
			}
			if y1 >= bounds.Dy() {
				y1 = bounds.Dy() - 1
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

func init() {
	fx.DefaultRegistry.Register(NewScale(1, 1))
}
