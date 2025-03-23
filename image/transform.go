// image/transform.go
package image

import (
	"image"
	"math"
	"sync"
)

// TransformRotateScale applies a combined rotation and scaling transformation
// to the image around the specified center point. The transformation is defined as:
//
//	p_dest = center + factor * R(theta) * (p_src - center)
//
// This method computes the inverse mapping for each destination pixel.
// The original image dimensions are maintained, so parts of the transformed image that fall outside are clipped.
func (i *Image) TransformRotateScale(angle float64, factor float64, centerX, centerY int) *Image {
	i.Lock()
	defer i.Unlock()
	w := i.raw.Bounds().Dx()
	h := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	// Convert angle to radians.
	theta := angle * math.Pi / 180.0
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)

	// For each pixel in the destination image, compute the corresponding source pixel.
	for y := range h {
		for x := range w {
			// Compute offset from the transformation center.
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			// Inverse mapping: given destination p, find source p such that
			// p = center + factor * R(theta) * (p_src - center)
			// => p_src = center + R(-theta) * ((p - center) / factor)
			srcX := float64(centerX) + (dx*cosTheta+dy*sinTheta)/factor
			srcY := float64(centerY) + (-dx*sinTheta+dy*cosTheta)/factor
			sx := int(math.Round(srcX))
			sy := int(math.Round(srcY))
			if sx >= 0 && sx < w && sy >= 0 && sy < h {
				dst.Set(x, y, i.raw.At(sx, sy))
			}
		}
	}
	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}
