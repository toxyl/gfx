package parser

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/fs/net"
)

type Layer struct {
	data      *image.Image
	Source    string
	BlendMode string
	Alpha     float64
	Filter    *FX
}

func (l *Layer) String() string {
	filter := "               *"
	wfilter := len(filter)
	if l.Filter != nil {
		filter = fmt.Sprintf("%*s", wfilter, l.Filter.Name)
	}
	return fmt.Sprintf(
		"%16s %6.4f %s %s",
		l.BlendMode,
		l.Alpha,
		filter,
		l.Source,
	)
}

func (l *Layer) load() error {
	if l.Source != "" {
		if l.Source[0] == '$' {
			if i, err := strconv.Atoi(l.Source[1:]); err == nil {
				l.Source = flag.Arg(i)
				if l.Source == "" {
					return fmt.Errorf("missing argument $%d (hint: numbering starts at 0)", i)
				}
			}
		}

		// Try loading from URL first
		if net.IsURL(l.Source) {
			// Download the file from URL
			resp, err := http.Get(l.Source)
			if err != nil {
				return fmt.Errorf("failed to download from URL: %w", err)
			}
			defer resp.Body.Close()

			// Create a temporary file
			tempFile, err := os.CreateTemp("", "image-*.tmp")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()

			// Copy the response body to the temp file
			_, err = io.Copy(tempFile, resp.Body)
			if err != nil {
				return fmt.Errorf("failed to copy response to temp file: %w", err)
			}

			// Close the file to ensure all data is written
			tempFile.Close()

			// Load the image from the temp file
			data, err := image.LoadImage(tempFile.Name())
			if err == nil {
				l.data = data
				return nil
			}
			// If loading fails, continue to try other methods
		}

		// Try loading from file
		data, err := image.LoadImage(l.Source)
		if err != nil {
			return fmt.Errorf("failed to load image: %w", err)
		}
		l.data = data
	}
	return nil
}

func (l *Layer) Render(w, h int) (*image.Image, error) {
	if err := l.load(); err != nil {
		return nil, err
	}
	if l.data == nil {
		return nil, fmt.Errorf("no image data available")
	}

	// Resize the image
	res, err := l.data.Resize(w, h, image.ResizeBilinear)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}

	// Apply filters if any
	if l.Filter != nil {
		for _, filter := range l.Filter.Get() {
			if filter != nil {
				var filterErr error
				stdImg := res.ToStandard()
				newImg, filterErr := filter.Apply(stdImg)
				if filterErr != nil {
					return nil, fmt.Errorf("failed to apply filter: %w", filterErr)
				}
				res, err = image.FromImage(newImg)
				if err != nil {
					return nil, fmt.Errorf("failed to convert image: %w", err)
				}
			}
		}
	}
	return res, nil
}

func NewLayer(blendmode string, alpha float64, filter *FX, source string) *Layer {
	l := Layer{
		data:      nil,
		Source:    source,
		BlendMode: blendmode,
		Alpha:     alpha,
		Filter:    filter,
	}
	return &l
}
