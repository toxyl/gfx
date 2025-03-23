package image

import (
	"image"
	"sync"

	"github.com/toxyl/gfx/color/rgba"
)

func (i *Image) Resize(w, h int) *Image {
	i.Lock()
	defer i.Unlock()
	orgW, orgH := i.raw.Bounds().Max.X, i.raw.Bounds().Max.Y
	scaleX, scaleY := float64(orgW)/float64(w), float64(orgH)/float64(h)
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX, srcY := float64(x)*scaleX, float64(y)*scaleY
			newX, newY := int(srcX), int(srcY)
			r, g, b, a := 0.0, 0.0, 0.0, 0.0
			if newX >= 0 && newX < orgW && newY >= 0 && newY < orgH {
				rf, gf, bf, af := i.raw.At(newX, newY).RGBA()
				r = float64(rf)
				g = float64(gf)
				b = float64(bf)
				a = float64(af)
			}
			res.Set(x, y, rgba.New(uint8(r/257), uint8(g/257), uint8(b/257), uint8(a/257)).RGBA())
		}
	}
	return &Image{raw: res, path: i.path, mu: &sync.Mutex{}}
}

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
		i = i.Resize(w, h)
	}
	return i
}
