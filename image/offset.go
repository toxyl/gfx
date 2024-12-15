package image

import (
	"image"
	"sync"
)

func (i *Image) Offset(dx, dy int) *Image {
	i.Lock()
	defer i.Unlock()
	w, h := i.raw.Bounds().Max.X, i.raw.Bounds().Max.Y
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if dx+x > w || dy+y > h {
				continue
			}
			res.Set(dx+x, dy+y, i.GetRGBA(x, y).RGBA())
		}
	}
	return &Image{raw: res, mu: &sync.Mutex{}}
}
