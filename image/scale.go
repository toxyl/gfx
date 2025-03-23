// image/scale.go
package image

import (
	"image"
	"math"
	"sync"
)

// Scale scales the image by the given factor around the specified center point.
// The original image dimensions are maintained, so parts of the scaled image that fall outside are clipped.
func (i *Image) Scale(factor float64, centerX, centerY int) *Image {
	factor += 1
	i.Lock()
	defer i.Unlock()
	w := i.raw.Bounds().Dx()
	h := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	// For each pixel in the destination image, compute the corresponding source pixel.
	for y := range h {
		for x := range w {
			// Compute the offset from the scaling center.
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			// Apply the inverse scaling transformation.
			srcX := float64(centerX) + dx/factor
			srcY := float64(centerY) + dy/factor
			sx := int(math.Round(srcX))
			sy := int(math.Round(srcY))
			if sx >= 0 && sx < w && sy >= 0 && sy < h {
				dst.Set(x, y, i.raw.At(sx, sy))
			}
		}
	}
	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}
