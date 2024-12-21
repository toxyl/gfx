package parser

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/toxyl/gfx/color/blend"
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

func (l *Layer) load() {
	if l.Source != "" {
		if l.Source[0] == CHAR_CLI_ARG {
			if i, err := strconv.Atoi(l.Source[1:]); err == nil {
				l.Source = flag.Arg(i)
				if l.Source == "" {
					panic("missing argument $" + fmt.Sprint(i) + " (hint: numbering starts at 0)")
				}
			}
		}
		l.data = image.NewFromURL(l.Source)
		if l.data == nil {
			// this wasn't a URL, maybe it's a file
			l.data = image.NewFromFile(l.Source)
		}
	}
}

func (l *Layer) LoadFromImage(i *image.Image) *Layer {
	l.data = i
	return l
}

func (l *Layer) SetFilters(filters ...*ImageFilter) *Layer {
	l.Filter = &CompiledFilter{
		Name:    "",
		Filters: filters,
	}
	return l
}

func (l *Layer) SetBlendmode(mode string) *Layer {
	l.BlendMode = mode
	return l
}

func (l *Layer) SetAlpha(alpha float64) *Layer {
	l.Alpha = alpha
	return l
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
		res2 := image.New(w, h)
		res2.Draw(res.Resize(l.Resize.W, l.Resize.H), 0, 0, l.Resize.W, l.Resize.H, (w-l.Resize.W)/2, (h-l.Resize.H)/2, l.Resize.W, l.Resize.H, blend.NORMAL, 1)
		res = res2
	}
	if l.Crop != nil && l.Crop.W > 0 && l.Crop.H > 0 {
		res = res.Crop(l.Crop.X, l.Crop.Y, l.Crop.W, l.Crop.H, false)
	}
	if l.Offset != nil && (l.Offset.X != 0 || l.Offset.Y != 0) {
		res = res.Offset(l.Offset.X, l.Offset.Y)
	}
	return res
}

func NewLayer() *Layer {
	l := Layer{
		data:      nil,
		Source:    "",
		BlendMode: string(blend.NORMAL),
		Alpha:     1,
		Crop:      nil,
		Offset:    nil,
		Resize:    nil,
		Filter:    nil,
	}
	return &l
}
