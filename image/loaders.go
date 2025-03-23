package image

import (
	"image"
	"strings"

	"github.com/toxyl/errors"
	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/jpg"
	"github.com/toxyl/gfx/net"
	"github.com/toxyl/gfx/png"
)

func loadFromURL(url string) (image.Image, error) {
	imgData, err := net.Download(url)
	if err != nil {
		return nil, err
	}
	u := strings.ToLower(url)
	if strings.HasSuffix(u, ".png") {
		return png.FromBytes(imgData)
	}

	if strings.HasSuffix(u, ".jpg") || strings.HasSuffix(u, ".jpeg") {
		return jpg.FromBytes(imgData)
	}
	return nil, errors.Newf("unknown format: %s", url)
}

func loadFromFile(path string) (image.Image, error) {
	imgData := flo.File(path).AsBytes()
	if len(imgData) == 0 {
		return nil, errors.Newf("no data found at %s", path)
	}
	u := strings.ToLower(path)
	if strings.HasSuffix(u, ".png") {
		return png.FromBytes(imgData)
	}

	if strings.HasSuffix(u, ".jpg") || strings.HasSuffix(u, ".jpeg") {
		return jpg.FromBytes(imgData)
	}
	return nil, errors.Newf("unknown format: %s", path)
}

// loadFromBytes generates an image from byte data using the given type. Available types: png, jpg and jpeg
func loadFromBytes(typ string, data []byte) (image.Image, error) {
	u := strings.ToLower(typ)
	if u == "png" {
		return png.FromBytes(data)
	}
	if u == "jpg" || u == "jpeg" {
		return jpg.FromBytes(data)
	}
	return nil, errors.Newf("unknown image type: %s", typ)
}
