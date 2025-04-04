// Package image provides a comprehensive set of image manipulation and processing utilities.
//
// This package implements a modern, type-safe approach to image processing with support for:
// - Multiple color models via the core/color package
// - Advanced blending operations via the core/blendmodes package
// - Efficient image transformations and manipulations
// - Composable processing pipelines
//
// The package is designed to be extensible and performant, with a clean API
// that aligns with modern Go practices.
package image

import (
	"fmt"
	stdImage "image"
	"image/draw"
	"sync"

	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/math"
)

// Image represents a 2D grid of pixels, each with color information.
// It is the core type for image manipulation in this package.
type Image struct {
	// Width and height of the image in pixels
	width, height int

	// The underlying image data
	data *stdImage.RGBA

	// Metadata about the image
	meta *Metadata

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Metadata stores information about an image
type Metadata struct {
	// Source information (filename, origin, etc.)
	Source string

	// Additional metadata (key-value pairs)
	Properties map[string]interface{}
}

// New creates a new image with the specified dimensions.
// The image is initialized with transparent pixels.
func New(width, height int) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: width=%d, height=%d", width, height)
	}

	img := &Image{
		width:  width,
		height: height,
		data:   stdImage.NewRGBA(stdImage.Rect(0, 0, width, height)),
		meta: &Metadata{
			Properties: make(map[string]interface{}),
		},
	}

	return img, nil
}

// FromImage creates a new Image from a standard library image.
func FromImage(img stdImage.Image) (*Image, error) {
	if img == nil {
		return nil, fmt.Errorf("cannot create image from nil")
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	newImg, err := New(width, height)
	if err != nil {
		return nil, err
	}

	// Draw the source image onto our RGBA
	draw.Draw(newImg.data, newImg.data.Bounds(), img, bounds.Min, draw.Src)

	return newImg, nil
}

// Width returns the width of the image in pixels.
func (i *Image) Width() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.width
}

// Height returns the height of the image in pixels.
func (i *Image) Height() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.height
}

// Bounds returns the bounds of the image.
func (i *Image) Bounds() stdImage.Rectangle {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data.Bounds()
}

// Size returns the width and height of the image.
func (i *Image) Size() (width, height int) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.width, i.height
}

// SetPixel sets the color of a pixel at (x, y).
func (i *Image) SetPixel(x, y int, c *color.RGBA64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if x < 0 || x >= i.width || y < 0 || y >= i.height {
		return fmt.Errorf("pixel coordinates out of bounds: x=%d, y=%d", x, y)
	}

	i.data.SetRGBA64(x, y, c.To16bit())
	return nil
}

// GetPixel returns the color of a pixel at (x, y).
func (i *Image) GetPixel(x, y int) (*color.RGBA64, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if x < 0 || x >= i.width || y < 0 || y >= i.height {
		return nil, fmt.Errorf("pixel coordinates out of bounds: x=%d, y=%d", x, y)
	}

	r, g, b, a := i.data.RGBA64At(x, y).RGBA()
	return &color.RGBA64{
		R: float64(r) / 65535.0,
		G: float64(g) / 65535.0,
		B: float64(b) / 65535.0,
		A: float64(a) / 65535.0,
	}, nil
}

// ToStandard returns the standard library image.Image interface.
func (i *Image) ToStandard() stdImage.Image {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data
}

// Clone creates a deep copy of the image.
func (i *Image) Clone() (*Image, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	newImg, err := New(i.width, i.height)
	if err != nil {
		return nil, err
	}

	// Copy pixel data
	draw.Draw(newImg.data, newImg.data.Bounds(), i.data, i.data.Bounds().Min, draw.Src)

	// Copy metadata
	if i.meta != nil {
		newImg.meta = &Metadata{
			Source: i.meta.Source,
		}

		// Copy properties
		if i.meta.Properties != nil {
			newImg.meta.Properties = make(map[string]interface{})
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}
	}

	return newImg, nil
}

// Process applies a function to each pixel of the image.
// The function receives the x, y coordinates and current color,
// and should return a new color for the pixel.
func (i *Image) Process(processor func(x, y int, c *color.RGBA64) (*color.RGBA64, error)) error {
	// By default, use parallel processing for better performance
	return i.ProcessParallel(processor)
}

// ProcessSequential applies a function sequentially to each pixel of the image.
// The function receives the x, y coordinates and current color,
// and should return a new color for the pixel.
// This method is useful when order of processing matters.
func (i *Image) ProcessSequential(processor func(x, y int, c *color.RGBA64) (*color.RGBA64, error)) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	var lastErr error

	for y := 0; y < i.height; y++ {
		for x := 0; x < i.width; x++ {
			r, g, b, a := i.data.RGBA64At(x, y).RGBA()
			currentColor := &color.RGBA64{
				R: float64(r) / 65535.0,
				G: float64(g) / 65535.0,
				B: float64(b) / 65535.0,
				A: float64(a) / 65535.0,
			}

			newColor, err := processor(x, y, currentColor)
			if err != nil {
				lastErr = err
				continue
			}

			i.data.SetRGBA64(x, y, newColor.To16bit())
		}
	}

	return lastErr
}

// ProcessParallel applies a function to each pixel of the image in parallel.
// The function receives the x, y coordinates and current color,
// and should return a new color for the pixel.
func (i *Image) ProcessParallel(processor func(x, y int, c *color.RGBA64) (*color.RGBA64, error)) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Create a new image to store the result
	result := stdImage.NewRGBA(i.data.Bounds())

	var wg sync.WaitGroup
	var errMutex sync.Mutex
	var lastErr error

	// Process each row in parallel
	for y := 0; y < i.height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()

			for x := 0; x < i.width; x++ {
				r, g, b, a := i.data.RGBA64At(x, y).RGBA()
				currentColor := &color.RGBA64{
					R: float64(r) / 65535.0,
					G: float64(g) / 65535.0,
					B: float64(b) / 65535.0,
					A: float64(a) / 65535.0,
				}

				newColor, err := processor(x, y, currentColor)
				if err != nil {
					errMutex.Lock()
					lastErr = err
					errMutex.Unlock()
					continue
				}

				result.SetRGBA64(x, y, newColor.To16bit())
			}
		}(y)
	}

	wg.Wait()

	// Copy result back to the original image
	i.data = result

	return lastErr
}

// Blend blends this image with another using the specified blend mode.
func (i *Image) Blend(other *Image, mode *blendmodes.IBlendMode, alpha float64) error {
	if other == nil {
		return fmt.Errorf("cannot blend with nil image")
	}

	if mode == nil {
		return fmt.Errorf("blend mode cannot be nil")
	}

	// Get dimensions of both images
	srcW, srcH := other.Size()
	dstW, dstH := i.Size()

	// Find the intersection
	width := math.Min(srcW, dstW)
	height := math.Min(srcH, dstH)

	// Process each pixel in the intersection
	return i.Process(func(x, y int, dstColor *color.RGBA64) (*color.RGBA64, error) {
		if x >= width || y >= height {
			return dstColor, nil
		}

		srcColor, err := other.GetPixel(x, y)
		if err != nil {
			return dstColor, err
		}

		result, err := mode.Blend(dstColor, srcColor, alpha)
		if err != nil {
			return dstColor, err
		}

		return result, nil
	})
}
