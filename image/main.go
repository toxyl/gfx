package image

import (
	"image"
	"image/draw"
	"sync"
	"time"

	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/jpg"
	"github.com/toxyl/gfx/net"
	"github.com/toxyl/gfx/png"
)

func toRGBAImage(img image.Image) *image.RGBA {
	switch t := img.(type) {
	case *image.RGBA:
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, t.Bounds(), t, image.Point{}, draw.Src)
		return c
	case *image.NRGBA:
		// Convert NRGBA to RGBA
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, t.Bounds(), t, image.Point{}, draw.Src)
		return c
	case *image.RGBA64:
		// Convert RGBA64 to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				r, g, b, a := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(r>>8, g>>8, b>>8, a>>8).RGBA())
			}
		}
		return c
	case *image.NRGBA64:
		// Convert NRGBA64 to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				r, g, b, a := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(r>>8, g>>8, b>>8, a>>8).RGBA())
			}
		}
		return c
	case *image.CMYK:
		// Convert CMYK to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				r, g, b, a := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(r>>8, g>>8, b>>8, a>>8).RGBA())
			}
		}
		return c
	case *image.Alpha:
		// Convert Alpha to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				_, _, _, a := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(0, 0, 0, a>>8).RGBA())
			}
		}
		return c
	case *image.Alpha16:
		// Convert Alpha16 to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				_, _, _, a := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(0, 0, 0, a>>8).RGBA())
			}
		}
		return c
	case *image.Gray:
		// Convert Gray to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				gray, _, _, _ := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(gray>>8, gray>>8, gray>>8, 255).RGBA())
			}
		}
		return c
	case *image.Gray16:
		// Convert Gray16 to RGBA
		c := image.NewRGBA(t.Bounds())
		for y := 0; y < t.Bounds().Dy(); y++ {
			for x := 0; x < t.Bounds().Dx(); x++ {
				gray, _, _, _ := t.At(x, y).RGBA()
				c.Set(x, y, rgba.New(gray>>8, gray>>8, gray>>8, 255).RGBA())
			}
		}
		return c
	case *image.YCbCr:
		// Convert YCbCr to RGBA
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, t.Bounds(), t, image.Point{}, draw.Src)
		return c
	case *image.NYCbCrA:
		// Convert NYCbCrA to RGBA
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, t.Bounds(), t, image.Point{}, draw.Src)
		return c
	case *image.Paletted:
		// Convert Paletted to RGBA
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, t.Bounds(), t, image.Point{}, draw.Src)
		return c
	}
	return img.(*image.RGBA)
}

type Image struct {
	mu  *sync.Mutex
	raw *image.RGBA
}

func (i *Image) SaveAsPNG(path string) { png.Save(i.raw, path) }
func (i *Image) SaveAsJPG(path string) { jpg.Save(i.raw, path) }
func (i *Image) W() int                { return i.raw.Bounds().Dx() }
func (i *Image) H() int                { return i.raw.Bounds().Dy() }
func (i *Image) Lock()                 { i.mu.Lock() }
func (i *Image) Unlock()               { i.mu.Unlock() }

func (i *Image) Set(img *image.RGBA) {
	if img == nil {
		return
	}
	i.Lock()
	defer i.Unlock()
	i.raw = img
}

func (i *Image) Get() *image.RGBA {
	i.Lock()
	defer i.Unlock()
	return i.raw
}

func New(w, h int) *Image {
	i := &Image{raw: image.NewRGBA(image.Rect(0, 0, w, h)), mu: &sync.Mutex{}}
	return i.FillRGBA(0, 0, w, h, rgba.New(0, 0, 0, 0))
}

func NewWithColor(w, h int, col rgba.RGBA) *Image {
	i := &Image{raw: image.NewRGBA(image.Rect(0, 0, w, h)), mu: &sync.Mutex{}}
	return i.FillRGBA(0, 0, w, h, &col)
}

func NewFromURL(url string) *Image {
	if !net.IsURL(url) {
		return nil
	}

	const maxAttempts = 3
	delay := time.Second

	var img image.Image
	var err error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		img, err = loadFromURL(url)
		if err == nil {
			return &Image{raw: toRGBAImage(img), mu: &sync.Mutex{}}
		}

		if attempt < maxAttempts-1 {
			time.Sleep(delay)
			delay *= 2
		}
	}

	return nil
}

func NewFromFile(path string) *Image {
	if i, err := loadFromFile(path); err == nil {
		return &Image{raw: toRGBAImage(i), mu: &sync.Mutex{}}
	}
	return nil
}

func NewFromBytes(typ string, b []byte) *Image {
	if i, err := loadFromBytes(typ, b); err == nil {
		return &Image{raw: toRGBAImage(i), mu: &sync.Mutex{}}
	}
	return nil
}

func NewFromImage(img image.Image) *Image {
	return &Image{raw: toRGBAImage(img), mu: &sync.Mutex{}}
}
