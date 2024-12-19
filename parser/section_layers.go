package parser

import (
	"strconv"
	"strings"
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
