package fx

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// CropFunction represents a function that crops an image
type CropFunction struct {
	*fx.BaseFunction
	x, y, width, height float64 // Crop rectangle coordinates and dimensions (0.0 to 1.0)
}

// Function arguments
var cropArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "crop adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewCrop(x, y, width, height float64) *CropFunction {
	return &CropFunction{
		BaseFunction: fx.NewBaseFunction("crop", "Applies crop transformation to an image", color.New(0, 0, 0, 1), cropArgs),
		x:            math.Clamp(x, 0.0, 1.0),
		y:            math.Clamp(y, 0.0, 1.0),
		width:        math.Clamp(width, 0.0, 1.0),
		height:       math.Clamp(height, 0.0, 1.0),
	}
}

// Apply implements the Function interface
func (f *CropFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Calculate crop rectangle in pixels
	cropX := int(float64(bounds.Dx()) * f.x)
	cropY := int(float64(bounds.Dy()) * f.y)
	cropWidth := int(float64(bounds.Dx()) * f.width)
	cropHeight := int(float64(bounds.Dy()) * f.height)

	// Ensure crop rectangle is within bounds
	cropX = math.Max(cropX, bounds.Min.X)
	cropY = math.Max(cropY, bounds.Min.Y)
	cropWidth = math.Min(cropWidth, bounds.Max.X-cropX)
	cropHeight = math.Min(cropHeight, bounds.Max.Y-cropY)

	// Create new image with cropped bounds
	cropBounds := image.Rect(0, 0, cropWidth, cropHeight)
	dst := image.NewRGBA(cropBounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Copy pixels from source to destination
	for y := 0; y < cropHeight; y++ {
		for x := 0; x < cropWidth; x++ {
			srcX := cropX + x
			srcY := cropY + y

			// Skip if source coordinates are out of bounds
			if srcX < bounds.Min.X || srcX >= bounds.Max.X ||
				srcY < bounds.Min.Y || srcY >= bounds.Max.Y {
				continue
			}

			// Copy pixel
			dst.Set(x, y, temp.At(srcX, srcY))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewCrop(0, 0, 1, 1))
}

func Apply(img *image.Image, left, right, top, bottom float64) *image.Image {
	left = math.Clamp(left, 0, 1)
	right = math.Clamp(right, 0, 1)
	top = math.Clamp(top, 0, 1)
	bottom = math.Clamp(bottom, 0, 1)

	bounds := (*img).Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	x := int(left * w)
	y := int(top * h)
	x2 := int((1 - right) * w)
	y2 := int((1 - bottom) * h)

	// Create new image with cropped bounds
	cropBounds := image.Rect(0, 0, x2-x, y2-y)
	dst := image.NewRGBA(cropBounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Copy pixels from source to destination
	for dy := 0; dy < y2-y; dy++ {
		for dx := 0; dx < x2-x; dx++ {
			srcX := x + dx
			srcY := y + dy
			dst.Set(dx, dy, temp.At(srcX, srcY))
		}
	}

	*img = dst
	return img
}
