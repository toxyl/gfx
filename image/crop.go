package image

import (
	"image"
	"image/color"
	"sync"

	coreimage "github.com/toxyl/gfx/core/image"
)

// RGBA64 represents a color with red, green, blue and alpha channels.
type RGBA64 struct {
	R float64
	G float64
	B float64
	A float64
}

// To8bit converts the RGBA64 to an 8-bit RGBA color.
func (c RGBA64) To8bit() color.RGBA {
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}

// Crop returns a cropped copy of the image.
// This function provides backward compatibility with the old API.
func (i *Image) Crop(x, y, w, h int) *Image {
	coreImg, err := coreimage.FromImage(i.raw)
	if err != nil {
		return i // Return original on error
	}

	cropped, err := coreImg.Crop(x, y, w, h)
	if err != nil {
		return i // Return original on error
	}

	result := &Image{
		raw:  cropped.ToStandard().(*image.RGBA),
		path: i.path,
		mu:   i.mu,
	}

	return result
}

// CropAround crops the image around a center point.
// This function provides backward compatibility with the old API.
func (i *Image) CropAround(cx, cy, w, h int) *Image {
	return i.Crop(cx-w/2, cy-h/2, w, h)
}

// CropCenter crops the image to the specified dimensions, keeping the center point of the original image.
// This function provides backward compatibility with the old API.
func (i *Image) CropCenter(w, h int) *Image {
	iw, ih := i.W(), i.H()
	return i.Crop((iw-w)/2, (ih-h)/2, w, h)
}

// AutoCrop automatically crops the image to remove borders of the specified color.
// This function provides backward compatibility with the old API.
func (i *Image) AutoCrop(bgColor color.Color, tolerance float64) *Image {
	bounds := i.raw.Bounds()
	x1, y1, x2, y2 := bounds.Min.X, bounds.Min.Y, bounds.Max.X-1, bounds.Max.Y-1
	r0, g0, b0, _ := bgColor.RGBA()

	// Scan from left
	for x := x1; x <= x2; x++ {
		found := false
		for y := y1; y <= y2; y++ {
			r, g, b, _ := i.raw.At(x, y).RGBA()
			if !isColorEqual(r, g, b, r0, g0, b0, tolerance) {
				found = true
				break
			}
		}
		if found {
			x1 = x
			break
		}
	}

	// Scan from right
	for x := x2; x >= x1; x-- {
		found := false
		for y := y1; y <= y2; y++ {
			r, g, b, _ := i.raw.At(x, y).RGBA()
			if !isColorEqual(r, g, b, r0, g0, b0, tolerance) {
				found = true
				break
			}
		}
		if found {
			x2 = x
			break
		}
	}

	// Scan from top
	for y := y1; y <= y2; y++ {
		found := false
		for x := x1; x <= x2; x++ {
			r, g, b, _ := i.raw.At(x, y).RGBA()
			if !isColorEqual(r, g, b, r0, g0, b0, tolerance) {
				found = true
				break
			}
		}
		if found {
			y1 = y
			break
		}
	}

	// Scan from bottom
	for y := y2; y >= y1; y-- {
		found := false
		for x := x1; x <= x2; x++ {
			r, g, b, _ := i.raw.At(x, y).RGBA()
			if !isColorEqual(r, g, b, r0, g0, b0, tolerance) {
				found = true
				break
			}
		}
		if found {
			y2 = y
			break
		}
	}

	width := x2 - x1 + 1
	height := y2 - y1 + 1
	return i.Crop(x1, y1, width, height)
}

// isColorEqual checks if two colors are approximately equal within a tolerance.
func isColorEqual(r1, g1, b1, r2, g2, b2 uint32, tolerance float64) bool {
	rd := float64(r1) - float64(r2)
	gd := float64(g1) - float64(g2)
	bd := float64(b1) - float64(b2)
	diff := (rd*rd + gd*gd + bd*bd) / (65535.0 * 65535.0)
	return diff <= tolerance*tolerance
}

// CropCircle creates a clone of the image and crops it to a circular area defined by the given center and radius.
// When fullCrop is true, the returned image dimensions are reduced to the minimal bounding square of the circle (2*radius × 2*radius).
// When fullCrop is false, the returned image has the same dimensions as the original image, with pixels outside the circle set to transparent.
func (i *Image) CropCircle(centerX, centerY, radius int, fullCrop bool) *Image {
	i.Lock()
	defer i.Unlock()
	rgba := RGBA64{R: 0, G: 0, B: 0, A: 0}
	transp := rgba.To8bit()

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
