package image

import (
	"image"
	"sync"
)

func (i *Image) FlipHorizontal() *Image {
	i.Lock()
	defer i.Unlock()
	orgW, orgH := i.raw.Bounds().Max.X, i.raw.Bounds().Max.Y
	res := image.NewRGBA(image.Rect(0, 0, orgW, orgH))
	for y := range orgH {
		for x := range orgW {
			res.Set(orgW-1-x, y, i.raw.At(x, y)) // Mirror the x-coordinate.
		}
	}
	return &Image{raw: res, path: i.path, mu: &sync.Mutex{}}
}

func (i *Image) FlipVertical() *Image {
	i.Lock()
	defer i.Unlock()
	orgW, orgH := i.raw.Bounds().Max.X, i.raw.Bounds().Max.Y
	res := image.NewRGBA(image.Rect(0, 0, orgW, orgH))
	for y := range orgH {
		for x := range orgW {
			res.Set(x, orgH-1-y, i.raw.At(x, y)) // Mirror the y-coordinate.
		}
	}
	return &Image{raw: res, path: i.path, mu: &sync.Mutex{}}
}
