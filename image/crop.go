package image

import (
	"image"
	"sync"
)

// Crop creates a clone of the image and crops it according to the `fullCrop` parameter:
//
// - `true`: only returns the crop area
//
// - `false`: returns the entire image, everything outside the crop area will we set to transparent
func (i *Image) Crop(x, y, w, h int, fullCrop bool) *Image {
	i.Lock()
	defer i.Unlock()
	if fullCrop {
		// this is a "real" crop operation
		// where the output image
		// dimensions match the crop area
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for yDst := 0; yDst < h; yDst++ {
			for xDst := 0; xDst < w; xDst++ {
				dst.Set(xDst, yDst, i.raw.RGBAAt(x+xDst, y+yDst))
			}
		}
		return &Image{raw: dst, mu: &sync.Mutex{}}
	}

	// this is a crop that only changes pixels
	// outside the crop area to be transparent
	iw := i.W()
	ih := i.H()
	dst := image.NewRGBA(image.Rect(0, 0, iw, ih))
	for yDst := x; yDst < h; yDst++ {
		for xDst := y; xDst < w; xDst++ {
			dst.Set(xDst, yDst, i.raw.RGBAAt(xDst, yDst))
		}
	}
	return &Image{raw: dst, mu: &sync.Mutex{}}
}
