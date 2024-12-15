package parser

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/toxyl/gfx/image"
)

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseLayer(line string, filters map[string]*CompiledFilter) Layer {
	parts := strings.Fields(line)
	blendMode := parts[0]
	alpha, _ := strconv.ParseFloat(parts[1], 64)
	filterName := parts[2]
	var crop *Crop
	var offset *Offset
	var resize *Resize
	var src string

	for i := 3; i < len(parts); {
		switch parts[i] {
		case LAYER_CROP:
			crop = &Crop{
				X: parseInt(parts[i+1]),
				Y: parseInt(parts[i+2]),
				W: parseInt(parts[i+3]),
				H: parseInt(parts[i+4]),
			}
			i += 5
		case LAYER_OFFSET:
			offset = &Offset{
				X: parseInt(parts[i+1]),
				Y: parseInt(parts[i+2]),
			}
			i += 3
		case LAYER_RESIZE:
			resize = &Resize{
				W: parseInt(parts[i+1]),
				H: parseInt(parts[i+2]),
			}
			i += 3
		default:
			src = strings.Join(parts[i:], STR_SPACE)
			i = len(parts)
		}
	}

	return Layer{
		Source:    src,
		BlendMode: blendMode,
		Alpha:     alpha,
		Crop:      crop,
		Offset:    offset,
		Resize:    resize,
		Filter:    filters[filterName],
	}
}

func (cl *Layer) load() {
	if cl.Source != "" {
		if cl.Source[0] == CHAR_CLI_ARG {
			if i, err := strconv.Atoi(cl.Source[1:]); err == nil {
				cl.Source = flag.Arg(i)
				if cl.Source == "" {
					panic("missing argument $" + fmt.Sprint(i) + " (hint: numbering starts at 0)")
				}
			}
		}
		cl.data = image.NewFromURL(cl.Source)
		if cl.data == nil {
			// this wasn't a URL, maybe it's a file
			cl.data = image.NewFromFile(cl.Source)
		}
	}
}

func (l *Layer) Render(w, h int) *image.Image {
	l.load()
	res := l.data.Resize(w, h)
	if l.Filter != nil {
		for _, filter := range l.Filter.Get() {
			if filter != nil {
				res = filter.Apply(res)
			}
		}
	}
	if l.Resize != nil && l.Resize.W > 0 && l.Resize.H > 0 {
		res = res.Resize(l.Resize.W, l.Resize.H)
	}
	if l.Crop != nil && l.Crop.W > 0 && l.Crop.H > 0 {
		res = res.Crop(l.Crop.X, l.Crop.Y, l.Crop.W, l.Crop.H, false)
	}
	if l.Offset != nil && l.Offset.X > 0 && l.Offset.Y > 0 {
		res = res.Offset(l.Offset.X, l.Offset.Y)
	}
	return res
}
