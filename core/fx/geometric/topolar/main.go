package fx

import (
	"image"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

const degToRadMultiplier = math.Pi / 180.0

// TopolarFunction represents a function that converts rectangular coordinates to polar coordinates
type TopolarFunction struct {
	*fx.BaseFunction
	angleStart float64
	angleEnd   float64
	rotation   float64
	fisheye    float64
}

// NewTopolar creates a new polar coordinate conversion function
// Function arguments
var topolarArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "topolar adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewTopolar(angleStart, angleEnd, rotation, fisheye float64) *TopolarFunction {
	return &TopolarFunction{
		BaseFunction: fx.NewBaseFunction("to-polar", "Applies topolar transformation to an image", color.New(0, 0, 0, 1), topolarArgs),
		angleStart:   math.Clamp(angleStart, 0, 360),
		angleEnd:     math.Clamp(angleEnd, 0, 360),
		rotation:     math.Clamp(rotation, -360, 360),
		fisheye:      math.Clamp(fisheye, -1, 1),
	}
}

// Apply implements the Function interface
func (f *TopolarFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	maxR := math.Min(cx, cy)

	// Precompute rotation parameters
	angleRad := f.rotation * degToRadMultiplier
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)

	// For each pixel in the final output image:
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Apply inverse rotation to get the corresponding coordinate in the pre-rotated polar image
			xt := float64(x) - cx
			yt := float64(y) - cy
			px := cosA*xt + sinA*yt + cx
			py := -sinA*xt + cosA*yt + cy

			// Compute vector from the center
			dx := px - cx
			dy := py - cy
			r := math.Sqrt(dx*dx + dy*dy)

			// Apply fisheye effect by remapping the normalized radius
			norm := r / maxR
			if f.fisheye != 0 {
				norm = math.Pow(norm, 1.0/(1.0+f.fisheye))
				r = norm * maxR
			}

			// Compute angle (in degrees) relative to the center and adjust by angleStart
			theta := (math.Atan2(dy, dx) * 180.0 / math.Pi) - f.angleStart
			if theta < 0 {
				theta += 360
			}

			// Determine angular proportion
			var proportion float64
			if f.angleEnd >= f.angleStart {
				if theta < f.angleStart || theta > f.angleEnd {
					continue
				}
				proportion = (theta - f.angleStart) / (f.angleEnd - f.angleStart)
			} else {
				if theta < f.angleStart && theta > f.angleEnd {
					continue
				}
				totalRange := (360 - f.angleStart) + f.angleEnd
				if theta >= f.angleStart {
					proportion = (theta - f.angleStart) / totalRange
				} else {
					proportion = (theta + (360 - f.angleStart)) / totalRange
				}
			}

			// Map the radial distance to the vertical coordinate in the source image
			srcY := int(r / maxR * float64(h-1))
			if srcY < 0 {
				srcY = 0
			} else if srcY >= h {
				srcY = h - 1
			}

			// Map the angular proportion to the horizontal coordinate in the source image
			srcX := int(proportion * float64(w-1))
			if srcX < 0 {
				srcX = 0
			} else if srcX >= w {
				srcX = w - 1
			}

			// Set the destination pixel using normalized RGBA64 values
			srcColor := img.At(srcX, srcY)
			dst.Set(x, y, srcColor)
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewTopolar(0, 360, 0, 0))
}
