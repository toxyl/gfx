package image

import (
	"image"
	"sync"
)

func (i *Image) Crop(x, y, w, h int) *Image {
	i.Lock()
	defer i.Unlock()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for yDst := 0; yDst < h; yDst++ {
		for xDst := 0; xDst < w; xDst++ {
			dst.Set(xDst, yDst, i.raw.RGBAAt(x+xDst, y+yDst))
		}
	}
	return &Image{raw: dst, mu: &sync.Mutex{}}
}
