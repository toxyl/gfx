package image

import (
	"fmt"
	stdImage "image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// LoadImage loads an image from a file.
// The image format is determined by the file extension.
func LoadImage(filename string) (*Image, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := stdImage.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	// Convert to our Image type
	result, err := FromImage(img)
	if err != nil {
		return nil, err
	}

	// Set the source metadata
	result.meta.Source = filename

	return result, nil
}

// SaveFormat represents the image format for saving.
type SaveFormat string

const (
	// PNG format with lossless compression (supports transparency)
	FormatPNG SaveFormat = "png"

	// JPEG format with lossy compression (does not support transparency)
	FormatJPEG SaveFormat = "jpeg"
)

// SaveOptions contains options for saving an image.
type SaveOptions struct {
	// Format to save in (PNG or JPEG)
	Format SaveFormat

	// Quality for JPEG compression (0-100, higher is better quality)
	// Ignored for PNG format
	Quality int
}

// DefaultSaveOptions returns default save options (PNG with maximum quality).
func DefaultSaveOptions() *SaveOptions {
	return &SaveOptions{
		Format:  FormatPNG,
		Quality: 100,
	}
}

// Save saves the image to a file with the specified options.
// If options is nil, default options are used.
func (i *Image) Save(filename string, options *SaveOptions) error {
	if options == nil {
		options = DefaultSaveOptions()
	}

	// Make sure the directory exists
	dir := filepath.Dir(filename)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Lock the image for reading
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Determine format from options or filename extension if not specified
	format := options.Format
	if format == "" {
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".png":
			format = FormatPNG
		case ".jpg", ".jpeg":
			format = FormatJPEG
		default:
			return fmt.Errorf("unsupported file format: %s", ext)
		}
	}

	// Encode the image
	switch format {
	case FormatPNG:
		if err := png.Encode(file, i.data); err != nil {
			return fmt.Errorf("error encoding PNG: %w", err)
		}

	case FormatJPEG:
		// JPEG doesn't support alpha, so we need to flatten the image on a white background
		if hasTransparency(i.data) {
			// Create a new RGBA image with white background
			bounds := i.data.Bounds()
			bg := stdImage.NewRGBA(bounds)

			// Fill with white
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					bg.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				}
			}

			// Composite the original image on top
			draw.Draw(bg, bounds, i.data, bounds.Min, draw.Over)

			// Encode the flattened image
			quality := options.Quality
			if quality <= 0 || quality > 100 {
				quality = 90 // Default quality
			}

			if err := jpeg.Encode(file, bg, &jpeg.Options{Quality: quality}); err != nil {
				return fmt.Errorf("error encoding JPEG: %w", err)
			}
		} else {
			// Image has no transparency, encode directly
			quality := options.Quality
			if quality <= 0 || quality > 100 {
				quality = 90 // Default quality
			}

			if err := jpeg.Encode(file, i.data, &jpeg.Options{Quality: quality}); err != nil {
				return fmt.Errorf("error encoding JPEG: %w", err)
			}
		}

	default:
		return fmt.Errorf("unsupported save format: %s", format)
	}

	return nil
}

// SavePNG saves the image as a PNG file.
func (i *Image) SavePNG(filename string) error {
	return i.Save(filename, &SaveOptions{Format: FormatPNG})
}

// SaveJPEG saves the image as a JPEG file with the specified quality (0-100).
func (i *Image) SaveJPEG(filename string, quality int) error {
	return i.Save(filename, &SaveOptions{Format: FormatJPEG, Quality: quality})
}

// hasTransparency checks if an image has any transparent or semi-transparent pixels.
func hasTransparency(img *stdImage.RGBA) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 65535 { // Not fully opaque
				return true
			}
		}
	}
	return false
}
