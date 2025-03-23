package image

import (
	"image"
	"math"
	"sync"
)

// Rotate rotates the image by an arbitrary angle (in degrees) around a specified center.
// The output image has the same dimensions as the source image. Pixels falling outside the source bounds are discarded.
// Positive angles rotate the image clockwise.
func (i *Image) Rotate(angle float64, centerX, centerY float64) *Image {
	i.Lock()
	defer i.Unlock()
	orgW, orgH := i.raw.Bounds().Max.X, i.raw.Bounds().Max.Y

	// Convert the angle from degrees to radians.
	theta := angle * math.Pi / 180.0

	// The output image has the same dimensions as the original.
	res := image.NewRGBA(image.Rect(0, 0, orgW, orgH))

	// For each pixel in the destination image, compute the corresponding source coordinate using inverse rotation.
	for y := range orgH {
		for x := range orgW {
			// Compute offset from the rotation center.
			dx := float64(x) - centerX
			dy := float64(y) - centerY

			// Inverse rotation (rotate counterclockwise by theta) to find the corresponding source pixel.
			srcX := dx*math.Cos(theta) - dy*math.Sin(theta) + centerX
			srcY := dx*math.Sin(theta) + dy*math.Cos(theta) + centerY

			// Use nearest-neighbor sampling.
			srcXi := int(srcX + 0.5)
			srcYi := int(srcY + 0.5)
			if srcXi >= 0 && srcXi < orgW && srcYi >= 0 && srcYi < orgH {
				res.Set(x, y, i.raw.At(srcXi, srcYi))
			}
		}
	}
	return &Image{raw: res, path: i.path, mu: &sync.Mutex{}}
}
