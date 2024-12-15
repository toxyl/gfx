package parser

import (
	"fmt"

	"github.com/toxyl/gfx/image"
)

type Layer struct {
	data      *image.Image    `yaml:"-"`
	Source    string          `yaml:"src,omitempty"`
	BlendMode string          `yaml:"blend,omitempty"`
	Alpha     float64         `yaml:"alpha,omitempty"`
	Crop      *Crop           `yaml:"crop,omitempty"`
	Offset    *Offset         `yaml:"offset,omitempty"`
	Resize    *Resize         `yaml:"resize,omitempty"`
	Filter    *CompiledFilter `yaml:"filter,omitempty"`
}

func (l *Layer) String(compHasCrop, compHasResize, compHasOffset, compHasFilter bool) string {
	resize := "                "
	wresize := len(resize)
	if !compHasResize {
		resize = ""
		wresize = 0
	}
	if l.Resize != nil {
		resize = fmt.Sprintf("%*s", wresize, l.Resize.String())
	}
	crop := "                        "
	wcrop := len(crop)
	if !compHasCrop {
		crop = ""
		wcrop = 0
	}
	if l.Crop != nil {
		crop = fmt.Sprintf("%*s", wcrop, l.Crop.String())
	}
	offset := "                "
	woffset := len(offset)
	if !compHasOffset {
		offset = ""
		woffset = 0
	}
	if l.Offset != nil {
		offset = fmt.Sprintf("%*s", woffset, l.Offset.String())
	}
	filter := "               *"
	wfilter := len(filter)
	if !compHasFilter {
		filter = ""
		wfilter = 0
	}
	if l.Filter != nil {
		filter = l.Filter.Name
		filter = fmt.Sprintf("%*s", wfilter, l.Filter.Name)
	}
	return fmt.Sprintf(
		"%16s %6.4f %s %s %s %s %s",
		l.BlendMode,
		l.Alpha,
		filter,
		resize,
		crop,
		offset,
		l.Source,
	)
}
