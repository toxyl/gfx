package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// CropCircleFunction represents a function that crops an image to a circle
type CropCircleFunction struct {
	*fx.BaseFunction
	centerX, centerY float64 // Center coordinates (0.0 to 1.0)
	radius           float64 // Radius (0.0 to 1.0)
}

// Function arguments
var cropcircleArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "cropcircle adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewCropCircle creates a new circular crop function
func NewCropCircle(centerX, centerY, radius float64) *CropCircleFunction {
	return &CropCircleFunction{
		BaseFunction: fx.NewBaseFunction("cropcircle", "Applies cropcircle transformation to an image", color.New(0, 0, 0, 1), cropcircleArgs),
		centerX:      math.Max(0.0, math.Min(1.0, centerX)),
		centerY:      math.Max(0.0, math.Min(1.0, centerY)),
		radius:       math.Max(0.0, math.Min(1.0, radius)),
	}
}

// Apply implements the Function interface
func (f *CropCircleFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Calculate circle parameters in pixels
	centerX := int(float64(bounds.Dx()) * f.centerX)
	centerY := int(float64(bounds.Dy()) * f.centerY)
	radius := int(float64(math.Min(float64(bounds.Dx()), float64(bounds.Dy()))) * f.radius / 2.0)

	// Create new image with same bounds
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate distance from center
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// If pixel is outside circle, make it transparent
			if distance > float64(radius) {
				dst.Set(x, y, stdcolor.RGBA{0, 0, 0, 0})
				continue
			}

			// Copy pixel from source
			dst.Set(x, y, temp.At(x, y))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewCropCircle(0.5, 0.5, 1.0))
}

func Apply(img *image.Image, radius, offsetX, offsetY float64) (*image.Image, error) {
	radius = math.Max(0, math.Min(radius, math.MaxFloat64))
	offsetX = math.Max(-1, math.Min(1, offsetX))
	offsetY = math.Max(-1, math.Min(1, offsetY))

	bounds := (*img).Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())
	hw := w / 2
	hh := h / 2

	// Calculate center point with offset
	centerX := int(hw + offsetX*hw)
	centerY := int(hh + offsetY*hh)
	radiusPixels := int(radius * math.Max(w, h))

	// Create new image with same bounds
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate distance from center
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// If pixel is outside circle, make it transparent
			if distance > float64(radiusPixels) {
				dst.Set(x, y, stdcolor.RGBA{0, 0, 0, 0})
				continue
			}

			// Copy pixel from source
			dst.Set(x, y, temp.At(x, y))
		}
	}

	result := image.Image(dst)
	return &result, nil
}
