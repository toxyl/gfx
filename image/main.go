package image

import (
	"image"
	"image/draw"
	"path/filepath"
	"sync"

	"github.com/toxyl/gfx/core/color"
	coreimage "github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/fs/jpg"
	"github.com/toxyl/gfx/fs/net"
	"github.com/toxyl/gfx/fs/png"
)

// Image represents an image with RGBA color channels.
// This is a compatibility wrapper around the core/image package.
type Image struct {
	mu   *sync.Mutex
	path string
	raw  *image.RGBA
}

// Path returns the path of the image file.
func (i *Image) Path() string {
	return i.path
}

// SaveAsPNG saves the image as a PNG file.
func (i *Image) SaveAsPNG(path string) *Image {
	png.Save(i.raw, path)
	i.path = path
	return i
}

// SaveAsJPG saves the image as a JPG file.
func (i *Image) SaveAsJPG(path string) *Image {
	jpg.Save(i.raw, path)
	i.path = path
	return i
}

// W returns the width of the image.
func (i *Image) W() int {
	return i.raw.Bounds().Dx()
}

// H returns the height of the image.
func (i *Image) H() int {
	return i.raw.Bounds().Dy()
}

// CW returns the center X coordinate of the image.
func (i *Image) CW() int {
	return i.raw.Bounds().Dx() >> 1
}

// CH returns the center Y coordinate of the image.
func (i *Image) CH() int {
	return i.raw.Bounds().Dy() >> 1
}

// Lock locks the image for exclusive access.
func (i *Image) Lock() {
	i.mu.Lock()
}

// Unlock unlocks the image.
func (i *Image) Unlock() {
	i.mu.Unlock()
}

// Set sets the image to the specified RGBA image.
func (i *Image) Set(img *image.RGBA) {
	i.raw = img
}

// Get returns the RGBA image.
func (i *Image) Get() *image.RGBA {
	return i.raw
}

// New creates a new image with the specified dimensions.
func New(w, h int) *Image {
	return &Image{
		mu:  &sync.Mutex{},
		raw: image.NewRGBA(image.Rect(0, 0, w, h)),
	}
}

// NewWithColor creates a new image with the specified dimensions and color.
func NewWithColor(w, h int, col color.RGBA64) *Image {
	i := New(w, h)

	// Convert from core/color.RGBA64 to standard RGBA color
	rgba8 := col.To8bit()

	// Fill the image with the color
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i.raw.SetRGBA(x, y, rgba8)
		}
	}

	return i
}

// NewFromURL creates a new image from the specified URL.
func NewFromURL(url string) *Image {
	b, err := net.Download(url)
	if err != nil {
		return nil
	}
	return NewFromBytes(filepath.Ext(url), b)
}

// NewFromFile creates a new image from the specified file.
func NewFromFile(path string) *Image {
	img, err := png.FromFile(path)
	if err != nil {
		img, err = jpg.FromFile(path)
		if err != nil {
			return nil
		}
	}
	return &Image{
		mu:   &sync.Mutex{},
		path: path,
		raw:  toRGBAImage(img),
	}
}

// NewFromBytes creates a new image from the specified bytes.
func NewFromBytes(ext string, b []byte) *Image {
	var img image.Image
	var err error

	switch filepath.Ext(ext) {
	case ".png", ".PNG":
		img, err = png.FromBytes(b)
	case ".jpg", ".jpeg", ".JPG", ".JPEG":
		img, err = jpg.FromBytes(b)
	default:
		return nil
	}

	if err != nil {
		return nil
	}

	return &Image{
		mu:  &sync.Mutex{},
		raw: toRGBAImage(img),
	}
}

// FromCoreImage creates a new Image from a core/image.Image
func FromCoreImage(img *coreimage.Image) *Image {
	if img == nil {
		return nil
	}

	return &Image{
		mu:  &sync.Mutex{},
		raw: img.ToStandard().(*image.RGBA),
	}
}

// ToCoreImage converts the Image to a core/image.Image
func (i *Image) ToCoreImage() *coreimage.Image {
	if i == nil || i.raw == nil {
		return nil
	}

	coreImg, err := coreimage.FromImage(i.raw)
	if err != nil {
		return nil
	}

	return coreImg
}

func toRGBAImage(img image.Image) *image.RGBA {
	switch t := img.(type) {
	case *image.RGBA:
		c := image.NewRGBA(t.Bounds())
		draw.Draw(c, c.Bounds(), t, t.Bounds().Min, draw.Src)
		return c
	default:
		// Use core/image conversion for all other types
		coreImg, err := coreimage.FromImage(img)
		if err != nil {
			// Fallback to standard conversion if core/image fails
			c := image.NewRGBA(img.Bounds())
			draw.Draw(c, c.Bounds(), img, img.Bounds().Min, draw.Src)
			return c
		}
		return coreImg.ToStandard().(*image.RGBA)
	}
}
