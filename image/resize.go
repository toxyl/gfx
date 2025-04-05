package image

import (
	"image"
	"image/color"

	coreimage "github.com/toxyl/gfx/core/image"
)

// These constants provide backward compatibility with old resize methods
const (
	ResizeTypeBilinear = iota
	ResizeTypeBicubic
	ResizeTypeNearest
)

// Resize scales the image to the new width and height using the specified algorithm.
// This function provides backward compatibility with the old API.
func (i *Image) Resize(nw, nh int, alg int) *Image {
	coreImg, err := coreimage.FromImage(i.raw)
	if err != nil {
		return i // Return original on error
	}

	// Map the old algorithm constants to the new ones
	var resizeMethod coreimage.ResizeMethod
	switch alg {
	case ResizeTypeBilinear:
		resizeMethod = coreimage.ResizeBilinear
	case ResizeTypeBicubic:
		resizeMethod = coreimage.ResizeBicubic
	case ResizeTypeNearest:
		resizeMethod = coreimage.ResizeNearest
	default:
		resizeMethod = coreimage.ResizeBilinear
	}

	resized, err := coreImg.Resize(nw, nh, resizeMethod)
	if err != nil {
		return i // Return original on error
	}

	result := &Image{
		raw:  resized.ToStandard().(*image.RGBA),
		path: i.path,
		mu:   i.mu,
	}

	return result
}

// ResizeRatio scales the image to a new size based on the given ratio.
// This function provides backward compatibility with the old API.
func (i *Image) ResizeRatio(ratio float64, alg int) *Image {
	w, h := i.W(), i.H()
	nw := int(float64(w) * ratio)
	nh := int(float64(h) * ratio)
	if nw <= 0 || nh <= 0 {
		return i
	}
	return i.Resize(nw, nh, alg)
}

// ResizeMax scales down the image to fit within the given maximum dimensions while maintaining the aspect ratio.
// This function provides backward compatibility with the old API.
func (i *Image) ResizeMax(maxW, maxH int, alg int) *Image {
	w, h := i.W(), i.H()
	if w <= maxW && h <= maxH {
		return i // Image already fits within the constraints
	}

	ratioW := float64(maxW) / float64(w)
	ratioH := float64(maxH) / float64(h)
	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}

	return i.ResizeRatio(ratio, alg)
}

// ResizeToMaxMP resizes the image so that its megapixel count doesn't exceed the specified value
// This function provides backward compatibility with the old API.
func (i *Image) ResizeToMaxMP(mpMax int) *Image {
	w := i.W()
	h := i.H()
	ow := w
	oh := h
	mp := w * h
	for mp > mpMax {
		w >>= 1
		h >>= 1
		mp = w * h
	}
	if ow != w || oh != h {
		i = i.Resize(w, h, ResizeTypeBilinear)
	}
	return i
}

// localRGBA represents a color with red, green, blue and alpha channels.
type localRGBA struct {
	R, G, B, A float64
}

// resizeBilinear resizes the image using bilinear interpolation.
func resizeBilinear(src image.Image, dst *image.RGBA, w, h, nw, nh int) {
	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			// Calculate the source coordinates
			sx := float64(x) * float64(w) / float64(nw)
			sy := float64(y) * float64(h) / float64(nh)

			// Get the four surrounding pixels
			x0, y0 := int(sx), int(sy)
			x1, y1 := min(x0+1, w-1), min(y0+1, h-1)

			// Calculate the fractional parts for interpolation
			fx := sx - float64(x0)
			fy := sy - float64(y0)

			// Get the color values of the four corners
			c00 := rgbaToFloat64(src.At(x0, y0))
			c01 := rgbaToFloat64(src.At(x0, y1))
			c10 := rgbaToFloat64(src.At(x1, y0))
			c11 := rgbaToFloat64(src.At(x1, y1))

			// Interpolate in the x direction
			c0 := interpolateColor(c00, c10, fx)
			c1 := interpolateColor(c01, c11, fx)

			// Interpolate in the y direction
			c := interpolateColor(c0, c1, fy)

			// Set the destination pixel
			dst.Set(x, y, float64ToRGBA(c))
		}
	}
}

// resizeBicubic resizes the image using bicubic interpolation.
func resizeBicubic(src image.Image, dst *image.RGBA, w, h, nw, nh int) {
	// Simple implementation - this could be improved
	resizeBilinear(src, dst, w, h, nw, nh)
}

// resizeNearest resizes the image using nearest neighbor interpolation.
func resizeNearest(src image.Image, dst *image.RGBA, w, h, nw, nh int) {
	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			// Calculate the source coordinates
			sx := int(float64(x) * float64(w) / float64(nw))
			sy := int(float64(y) * float64(h) / float64(nh))

			// Get the color at the source position
			c := src.At(sx, sy)

			// Set the destination pixel
			dst.Set(x, y, c)
		}
	}
}

// Helper function to convert image.Color to localRGBA
func rgbaToFloat64(c color.Color) localRGBA {
	r, g, b, a := c.RGBA()
	return localRGBA{
		R: float64(r) / 65535.0,
		G: float64(g) / 65535.0,
		B: float64(b) / 65535.0,
		A: float64(a) / 65535.0,
	}
}

// Helper function to convert localRGBA to color.Color
func float64ToRGBA(c localRGBA) color.Color {
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}

// Interpolate between two RGBA colors
func interpolateColor(c1, c2 localRGBA, t float64) localRGBA {
	return localRGBA{
		R: lerp(c1.R, c2.R, t),
		G: lerp(c1.G, c2.G, t),
		B: lerp(c1.B, c2.B, t),
		A: lerp(c1.A, c2.A, t),
	}
}

// Linear interpolation
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
