package image

import (
	"image"
	"image/draw"
	"sync"
)

func (i *Image) Clone() *Image {
	i.Lock()
	defer i.Unlock()
	rgbaImage := image.NewRGBA(i.raw.Bounds())
	draw.Draw(rgbaImage, rgbaImage.Bounds(), i.raw, image.Point{}, draw.Src)
	return &Image{raw: rgbaImage, mu: &sync.Mutex{}}
}
