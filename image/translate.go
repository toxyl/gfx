// image/translate.go
package image

import (
	"image"
	"sync"
)

// Translate translates the image by the specified absolute pixel offsets (x, y).
// The original image dimensions are maintained.
// If wrap is true, pixels that exit one side reappear on the opposite side (wrap-around).
// Otherwise, areas outside the original bounds are filled with transparent pixels.
func (i *Image) Translate(x, y int, wrap bool) *Image {
	i.Lock()
	defer i.Unlock()
	w := i.raw.Bounds().Dx()
	h := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for j := range h {
		for k := range w {
			srcX := k - x
			srcY := j - y
			if wrap {
				srcX = (srcX%w + w) % w
				srcY = (srcY%h + h) % h
				dst.Set(k, j, i.raw.At(srcX, srcY))
			} else if srcX >= 0 && srcX < w && srcY >= 0 && srcY < h {
				dst.Set(k, j, i.raw.At(srcX, srcY))
			}
		}
	}
	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}
