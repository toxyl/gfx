package image

import (
	"image"
	"sync"

	"github.com/toxyl/gfx/color/rgba"
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
		for y1 := range h {
			for x1 := range w {
				dst.Set(x1, y1, i.raw.RGBAAt(x+x1, y+y1))
			}
		}
		return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
	}

	// this is a crop that only changes pixels
	// outside the crop area to be transparent
	iw := i.W()
	ih := i.H()
	dst := image.NewRGBA(image.Rect(0, 0, iw, ih))
	for y1 := range h {
		for x1 := range w {
			dst.Set(x+x1, y+y1, i.raw.RGBAAt(x+x1, y+y1))
		}
	}
	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}

// CropCircle creates a clone of the image and crops it to a circular area defined by the given center and radius.
// When fullCrop is true, the returned image dimensions are reduced to the minimal bounding square of the circle (2*radius × 2*radius).
// When fullCrop is false, the returned image has the same dimensions as the original image, with pixels outside the circle set to transparent.
func (i *Image) CropCircle(centerX, centerY, radius int, fullCrop bool) *Image {
	i.Lock()
	defer i.Unlock()
	transp := rgba.New(0, 0, 0, 0).RGBA()

	if fullCrop {
		// Define the bounding box for the circle.
		diameter := 2 * radius
		dst := image.NewRGBA(image.Rect(0, 0, diameter, diameter))
		// The top-left of the new image corresponds to (centerX - radius, centerY - radius) in the original.
		startX := centerX - radius
		startY := centerY - radius

		// For each pixel in the destination image...
		for y := range diameter {
			for x := range diameter {
				// Determine if the pixel is inside the circle.
				dx := x - radius
				dy := y - radius
				if dx*dx+dy*dy <= radius*radius {
					origX := startX + x
					origY := startY + y
					// Only copy if the corresponding coordinate is within the original image.
					if origX >= 0 && origX < i.raw.Bounds().Dx() && origY >= 0 && origY < i.raw.Bounds().Dy() {
						dst.Set(x, y, i.raw.At(origX, origY))
					} else {
						// Out-of-bounds: set transparent.
						dst.Set(x, y, transp)
					}
				}
			}
		}
		return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
	}

	// For non-full crop: create an output image with the original dimensions,
	// but set pixels outside the circle to transparent.
	iw := i.raw.Bounds().Dx()
	ih := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, iw, ih))
	for y := range ih {
		for x := range iw {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				dst.Set(x, y, i.raw.At(x, y))
			}
		}
	}
	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}
