package image

import (
	"fmt"
	stdImage "image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/math"
)

// ResizeMethod represents different methods of image resizing
type ResizeMethod int

const (
	// ResizeNearest uses nearest neighbor interpolation (fastest, lowest quality)
	ResizeNearest ResizeMethod = iota

	// ResizeBilinear uses bilinear interpolation (good balance of speed and quality)
	ResizeBilinear

	// ResizeBicubic uses bicubic interpolation (better quality, slower)
	ResizeBicubic

	// ResizeLanczos uses Lanczos resampling (highest quality, slowest)
	ResizeLanczos
)

// Resize resizes the image to the specified dimensions using the given method.
func (i *Image) Resize(width, height int, method ResizeMethod) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid resize dimensions: width=%d, height=%d", width, height)
	}

	// Create a new image with the target dimensions
	newImg, err := New(width, height)
	if err != nil {
		return nil, err
	}

	// Lock the source image for reading
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Get dimensions of the source image
	srcW, srcH := i.width, i.height

	// Calculate scale factors
	xScale := float64(srcW) / float64(width)
	yScale := float64(srcH) / float64(height)

	// Process each pixel in the target image
	err = newImg.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
		// Map the target coordinates to source coordinates
		var c *color.RGBA64
		var err error

		switch method {
		case ResizeNearest:
			srcX := int(float64(x) * xScale)
			srcY := int(float64(y) * yScale)
			c, err = i.GetPixel(srcX, srcY)
			if err != nil {
				return nil, err
			}

		case ResizeBilinear:
			// Calculate source position with fraction
			srcX := float64(x) * xScale
			srcY := float64(y) * yScale

			// Get the four surrounding pixels
			x0, y0 := int(srcX), int(srcY)
			x1, y1 := math.Min(x0+1, srcW-1), math.Min(y0+1, srcH-1)

			// Calculate interpolation weights
			wx := srcX - float64(x0)
			wy := srcY - float64(y0)

			// Get the four corner colors
			c00, err := i.GetPixel(x0, y0)
			if err != nil {
				return nil, err
			}

			c10, err := i.GetPixel(x1, y0)
			if err != nil {
				return nil, err
			}

			c01, err := i.GetPixel(x0, y1)
			if err != nil {
				return nil, err
			}

			c11, err := i.GetPixel(x1, y1)
			if err != nil {
				return nil, err
			}

			// Bilinear interpolation
			r := lerp(
				lerp(c00.R, c10.R, wx),
				lerp(c01.R, c11.R, wx),
				wy,
			)

			g := lerp(
				lerp(c00.G, c10.G, wx),
				lerp(c01.G, c11.G, wx),
				wy,
			)

			b := lerp(
				lerp(c00.B, c10.B, wx),
				lerp(c01.B, c11.B, wx),
				wy,
			)

			a := lerp(
				lerp(c00.A, c10.A, wx),
				lerp(c01.A, c11.A, wx),
				wy,
			)

			// Create a new color using the calculated values
			newColor, err := color.NewRGBA64(r, g, b, a)
			if err != nil {
				return nil, err
			}
			c = newColor

		case ResizeBicubic, ResizeLanczos:
			// These are more complex and would be implemented with a proper image
			// processing library in a production environment.
			// For this example, we'll fall back to bilinear.
			return nil, fmt.Errorf("resize method not implemented: %v", method)
		}

		return c, nil
	})

	if err != nil {
		return nil, fmt.Errorf("resize operation failed: %w", err)
	}

	// Copy metadata
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add resize info to properties
		newImg.meta.Properties["resized_from_width"] = srcW
		newImg.meta.Properties["resized_from_height"] = srcH
		newImg.meta.Properties["resize_method"] = method
	}

	return newImg, nil
}

// Crop extracts a rectangular region from the image.
func (i *Image) Crop(x, y, width, height int) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid crop dimensions: width=%d, height=%d", width, height)
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	// Validate the crop region
	if x < 0 || y < 0 || x+width > i.width || y+height > i.height {
		return nil, fmt.Errorf("crop region outside image bounds: x=%d, y=%d, width=%d, height=%d", x, y, width, height)
	}

	// Create a new image for the cropped result
	newImg, err := New(width, height)
	if err != nil {
		return nil, err
	}

	// Define the source and destination rectangles
	srcRect := stdImage.Rect(x, y, x+width, y+height)
	dstRect := stdImage.Rect(0, 0, width, height)

	// Perform the crop
	draw.Draw(newImg.data, dstRect, i.data, srcRect.Min, draw.Src)

	// Copy metadata with crop info
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add crop info to properties
		newImg.meta.Properties["cropped_from_x"] = x
		newImg.meta.Properties["cropped_from_y"] = y
		newImg.meta.Properties["cropped_from_width"] = i.width
		newImg.meta.Properties["cropped_from_height"] = i.height
	}

	return newImg, nil
}

// Rotate rotates the image by the specified angle in degrees.
// This implementation handles arbitrary rotation angles.
func (i *Image) Rotate(angleDegrees float64) (*Image, error) {
	// Convert angle to radians
	angleRadians := angleDegrees * math.Pi / 180.0

	i.mu.RLock()
	defer i.mu.RUnlock()

	// Calculate dimensions of the rotated image
	srcW, srcH := float64(i.width), float64(i.height)

	// Calculate the dimensions of the rotated image
	cosAngle, sinAngle := math.Cos(angleRadians), math.Sin(angleRadians)
	newWidth := int(math.Ceil(math.Abs(srcW*cosAngle) + math.Abs(srcH*sinAngle)))
	newHeight := int(math.Ceil(math.Abs(srcW*sinAngle) + math.Abs(srcH*cosAngle)))

	// Create a new image with the calculated dimensions
	newImg, err := New(newWidth, newHeight)
	if err != nil {
		return nil, err
	}

	// Calculate the center of the original and new images
	srcCenterX, srcCenterY := srcW/2, srcH/2
	dstCenterX, dstCenterY := float64(newWidth)/2, float64(newHeight)/2

	// Process each pixel in the new image
	err = newImg.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
		// Calculate the corresponding position in the source image
		// by applying the inverse rotation
		dx, dy := float64(x)-dstCenterX, float64(y)-dstCenterY

		srcX := dx*cosAngle + dy*sinAngle + srcCenterX
		srcY := -dx*sinAngle + dy*cosAngle + srcCenterY

		// Check if the source pixel is within bounds
		if srcX < 0 || srcX >= srcW || srcY < 0 || srcY >= srcH {
			// Transparent pixel for out-of-bounds coordinates
			transparentColor, err := color.NewRGBA64(0, 0, 0, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to create transparent color: %w", err)
			}
			return transparentColor, nil
		}

		// Use bilinear interpolation for better quality
		x0, y0 := int(srcX), int(srcY)
		x1, y1 := math.Min(x0+1, int(srcW)-1), math.Min(y0+1, int(srcH)-1)

		wx, wy := srcX-float64(x0), srcY-float64(y0)

		c00, err := i.GetPixel(x0, y0)
		if err != nil {
			return nil, err
		}

		c10, err := i.GetPixel(x1, y0)
		if err != nil {
			return nil, err
		}

		c01, err := i.GetPixel(x0, y1)
		if err != nil {
			return nil, err
		}

		c11, err := i.GetPixel(x1, y1)
		if err != nil {
			return nil, err
		}

		r := lerp(lerp(c00.R, c10.R, wx), lerp(c01.R, c11.R, wx), wy)
		g := lerp(lerp(c00.G, c10.G, wx), lerp(c01.G, c11.G, wx), wy)
		b := lerp(lerp(c00.B, c10.B, wx), lerp(c01.B, c11.B, wx), wy)
		a := lerp(lerp(c00.A, c10.A, wx), lerp(c01.A, c11.A, wx), wy)

		resultColor, err := color.NewRGBA64(r, g, b, a)
		if err != nil {
			return nil, fmt.Errorf("failed to create result color: %w", err)
		}
		return resultColor, nil
	})

	if err != nil {
		return nil, fmt.Errorf("rotation operation failed: %w", err)
	}

	// Copy metadata with rotation info
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add rotation info to properties
		newImg.meta.Properties["rotated_angle_degrees"] = angleDegrees
		newImg.meta.Properties["original_width"] = i.width
		newImg.meta.Properties["original_height"] = i.height
	}

	return newImg, nil
}

// FlipHorizontal flips the image horizontally.
func (i *Image) FlipHorizontal() (*Image, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	newImg, err := New(i.width, i.height)
	if err != nil {
		return nil, err
	}

	// Process each pixel
	err = newImg.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
		// Horizontal flip: map x to (width - 1 - x)
		srcX := i.width - 1 - x
		return i.GetPixel(srcX, y)
	})

	if err != nil {
		return nil, fmt.Errorf("horizontal flip operation failed: %w", err)
	}

	// Copy metadata
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add flip info
		newImg.meta.Properties["flipped_horizontally"] = true
	}

	return newImg, nil
}

// FlipVertical flips the image vertically.
func (i *Image) FlipVertical() (*Image, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	newImg, err := New(i.width, i.height)
	if err != nil {
		return nil, err
	}

	// Process each pixel
	err = newImg.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
		// Vertical flip: map y to (height - 1 - y)
		srcY := i.height - 1 - y
		return i.GetPixel(x, srcY)
	})

	if err != nil {
		return nil, fmt.Errorf("vertical flip operation failed: %w", err)
	}

	// Copy metadata
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add flip info
		newImg.meta.Properties["flipped_vertically"] = true
	}

	return newImg, nil
}

// Translate moves the image by the specified offsets.
// Areas that are moved outside the image bounds are clipped,
// and new areas are filled with transparent pixels.
func (i *Image) Translate(xOffset, yOffset int) (*Image, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Create a new image with the same dimensions
	newImg, err := New(i.width, i.height)
	if err != nil {
		return nil, err
	}

	// Process each pixel in the new image
	err = newImg.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
		// Calculate source coordinates
		srcX, srcY := x-xOffset, y-yOffset

		// Check if the source pixel is within bounds
		if srcX < 0 || srcX >= i.width || srcY < 0 || srcY >= i.height {
			// Use transparent pixel for out-of-bounds
			transparentColor, err := color.NewRGBA64(0, 0, 0, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to create transparent color: %w", err)
			}
			return transparentColor, nil
		}

		// Copy the pixel from the source image
		return i.GetPixel(srcX, srcY)
	})

	if err != nil {
		return nil, fmt.Errorf("translate operation failed: %w", err)
	}

	// Copy metadata
	if i.meta != nil {
		newImg.meta.Source = i.meta.Source

		// Copy properties
		if i.meta.Properties != nil {
			for k, v := range i.meta.Properties {
				newImg.meta.Properties[k] = v
			}
		}

		// Add translation info
		newImg.meta.Properties["translated_x_offset"] = xOffset
		newImg.meta.Properties["translated_y_offset"] = yOffset
	}

	return newImg, nil
}

// lerp performs linear interpolation between a and b with weight t
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
