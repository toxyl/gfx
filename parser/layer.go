package parser

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/fs/net"
)

type Layer struct {
	data      *image.Image
	Source    string
	BlendMode string
	Alpha     float64
	Filter    *Filter
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
		if l.Source[0] == config.CHAR_CLI_ARG {
			if i, err := strconv.Atoi(l.Source[1:]); err == nil {
				l.Source = flag.Arg(i)
				if l.Source == "" {
					return fmt.Errorf("missing argument $%d (hint: numbering starts at 0)", i)
				}
			}
		}
		// Try loading from URL first
		if net.IsURL(l.Source) {
			data, err := image.LoadFromURL(l.Source)
			if err == nil {
				l.data = data
				return nil
			}
			// If URL loading fails, we'll continue to try loading as a file
		}

		// Try loading from file
		data, err := image.LoadFromFile(l.Source)
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
				res, filterErr = filter.Apply(res)
				if filterErr != nil {
					return nil, fmt.Errorf("failed to apply filter: %w", filterErr)
				}
			}
		}
	}
	return res, nil
}

func NewLayer(blendmode string, alpha float64, filter *Filter, source string) *Layer {
	l := Layer{
		data:      nil,
		Source:    source,
		BlendMode: blendmode,
		Alpha:     alpha,
		Filter:    filter,
	}
	return &l
}
