package parser

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/color/filter"
	"github.com/toxyl/gfx/filters/alphamap"
	"github.com/toxyl/gfx/filters/blur"
	"github.com/toxyl/gfx/filters/brightness"
	"github.com/toxyl/gfx/filters/colorshift"
	"github.com/toxyl/gfx/filters/contrast"
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/filters/edgedetect"
	"github.com/toxyl/gfx/filters/emboss"
	"github.com/toxyl/gfx/filters/enhance"
	"github.com/toxyl/gfx/filters/extract"
	"github.com/toxyl/gfx/filters/gamma"
	"github.com/toxyl/gfx/filters/gray"
	"github.com/toxyl/gfx/filters/hue"
	"github.com/toxyl/gfx/filters/huecontrast"
	"github.com/toxyl/gfx/filters/invert"
	"github.com/toxyl/gfx/filters/lum"
	"github.com/toxyl/gfx/filters/lumcontrast"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/filters/pastelize"
	"github.com/toxyl/gfx/filters/sat"
	"github.com/toxyl/gfx/filters/satcontrast"
	"github.com/toxyl/gfx/filters/sepia"
	"github.com/toxyl/gfx/filters/sharpen"
	"github.com/toxyl/gfx/filters/threshold"
	"github.com/toxyl/gfx/filters/vibrance"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

type Image = image.Image
type MetaData = meta.FilterMetaData
type Filter = ImageFilter
type FilterFn func(s *Filter, i *Image, m *MetaData)

type FilterMapEntry struct {
	Fn   FilterFn
	Meta *MetaData
}

func NewFilterMapEntry(meta *MetaData, fn FilterFn) *FilterMapEntry {
	fm := FilterMapEntry{
		Fn:   fn,
		Meta: meta,
	}
	return &fm
}

type FilterMap map[string]*FilterMapEntry

func (fm *FilterMap) Get(name string) (meta *MetaData, fn FilterFn) {
	if e, ok := (*fm)[name]; ok {
		return e.Meta, e.Fn
	}
	return nil, nil
}

func NewFilterMap(entries ...*FilterMapEntry) *FilterMap {
	fm := FilterMap{}
	for _, entry := range entries {
		fm[entry.Meta.Name] = entry
	}
	return &fm
}

var (
	Filters = NewFilterMap(
		NewFilterMapEntry(gray.Meta, func(s *Filter, i *Image, m *MetaData) {
			gray.Apply(i)
		}),
		NewFilterMapEntry(invert.Meta, func(s *Filter, i *Image, m *MetaData) {
			invert.Apply(i)
		}),
		NewFilterMapEntry(pastelize.Meta, func(s *Filter, i *Image, m *MetaData) {
			pastelize.Apply(i)
		}),
		NewFilterMapEntry(sepia.Meta, func(s *Filter, i *Image, m *MetaData) {
			sepia.Apply(i)
		}),
		NewFilterMapEntry(hue.Meta, func(s *Filter, i *Image, m *MetaData) {
			hue.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(sat.Meta, func(s *Filter, i *Image, m *MetaData) {
			sat.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(lum.Meta, func(s *Filter, i *Image, m *MetaData) {
			lum.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(huecontrast.Meta, func(s *Filter, i *Image, m *MetaData) {
			huecontrast.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(satcontrast.Meta, func(s *Filter, i *Image, m *MetaData) {
			satcontrast.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(lumcontrast.Meta, func(s *Filter, i *Image, m *MetaData) {
			lumcontrast.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(colorshift.Meta, func(s *Filter, i *Image, m *MetaData) {
			colorshift.Apply(i,
				s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)),
				s.GetOptionFloat64(m.NameOf(1), m.DefaultOf(1)),
				s.GetOptionFloat64(m.NameOf(2), m.DefaultOf(2)),
			)
		}),
		NewFilterMapEntry(brightness.Meta, func(s *Filter, i *Image, m *MetaData) {
			brightness.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(contrast.Meta, func(s *Filter, i *Image, m *MetaData) {
			contrast.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(gamma.Meta, func(s *Filter, i *Image, m *MetaData) {
			gamma.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(vibrance.Meta, func(s *Filter, i *Image, m *MetaData) {
			vibrance.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(enhance.Meta, func(s *Filter, i *Image, m *MetaData) {
			enhance.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(sharpen.Meta, func(s *Filter, i *Image, m *MetaData) {
			sharpen.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(blur.Meta, func(s *Filter, i *Image, m *MetaData) {
			blur.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(edgedetect.Meta, func(s *Filter, i *Image, m *MetaData) {
			edgedetect.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(emboss.Meta, func(s *Filter, i *Image, m *MetaData) {
			emboss.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(threshold.Meta, func(s *Filter, i *Image, m *MetaData) {
			threshold.Apply(i, s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)))
		}),
		NewFilterMapEntry(alphamap.Meta, func(s *Filter, i *Image, m *MetaData) {
			alphamap.Apply(i,
				s.GetOptionString(m.NameOf(0), m.DefaultOf(0)),
				s.GetOptionFloat64(m.NameOf(1), m.DefaultOf(1)),
				s.GetOptionFloat64(m.NameOf(2), m.DefaultOf(2)),
			)
		}),
		NewFilterMapEntry(extract.Meta, func(s *Filter, i *Image, m *MetaData) {
			extract.Apply(i, filter.ToColorFilter(
				s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)), s.GetOptionFloat64(m.NameOf(1), m.DefaultOf(1)), s.GetOptionFloat64(m.NameOf(2), m.DefaultOf(2)),
				s.GetOptionFloat64(m.NameOf(3), m.DefaultOf(3)), s.GetOptionFloat64(m.NameOf(4), m.DefaultOf(4)), s.GetOptionFloat64(m.NameOf(5), m.DefaultOf(5)),
				s.GetOptionFloat64(m.NameOf(6), m.DefaultOf(6)), s.GetOptionFloat64(m.NameOf(7), m.DefaultOf(7)), s.GetOptionFloat64(m.NameOf(8), m.DefaultOf(8)),
			))
		}),
		NewFilterMapEntry(convolution.Meta, func(s *Filter, i *Image, m *MetaData) {
			convolution.NewCustomFilter(
				s.GetOptionFloat64(m.NameOf(0), m.DefaultOf(0)),
				1.0+s.GetOptionFloat64(m.NameOf(1), m.DefaultOf(1)),
				s.GetOptionFloat64(m.NameOf(2), m.DefaultOf(2)),
				func(a float64) [][]float64 {
					return s.GetOptionMatrix(m.NameOf(3), m.DefaultOf(3).([][]float64))
				},
			).Apply(i)
		}),
	)
)

type ImageFilter struct {
	Type    string         `yaml:"type,omitempty"`
	Options map[string]any `yaml:"options,omitempty"`
}

func (f *ImageFilter) String(verbose bool) string {
	res := ""
	m, _ := Filters.Get(f.Type)
	for _, opt := range m.Args {
		k := opt.Name
		v := opt.Default
		if vopt, ok := f.Options[k]; ok {
			v = vopt
		}
		if verbose {
			k += STR_ASSIGN
		} else {
			k = ""
		}
		switch t := v.(type) {
		case string:
			res += STR_SPACE + k + STR_QUOTE + strings.ReplaceAll(t, STR_QUOTE, STR_ESCAPE+STR_QUOTE) + STR_QUOTE
		default:
			res += STR_SPACE + k + fmt.Sprint(t)
		}
	}
	return f.Type + STR_LPAREN + strings.TrimSpace(res) + STR_RPAREN
}

func (s *ImageFilter) GetOptionFloat64(option string, def any) float64 {
	v, ok := s.Options[option]
	if ok && v != nil {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		case int:
			return float64(val)
		case int8:
			return float64(val)
		case int16:
			return float64(val)
		case int32:
			return float64(val)
		case int64:
			return float64(val)
		case uint:
			return float64(val)
		case uint8:
			return float64(val)
		case uint16:
			return float64(val)
		case uint32:
			return float64(val)
		case uint64:
			return float64(val)
		default:
			// Do nothing, fall back to default value
		}
	}
	return def.(float64)
}

func (s *ImageFilter) GetOptionString(option string, def any) string {
	v, ok := s.Options[option]
	if ok && v != nil {
		return v.(string)
	}
	return def.(string)
}

func (s *ImageFilter) GetOptionMatrix(option string, def [][]float64) [][]float64 {
	v, ok := s.Options[option]
	if ok && v != nil {
		in := v.([]float64)

		// Calculate the number of rows and columns
		rows := int(math.Sqrt(float64(len(in))))

		// Check if the length of the input slice is a perfect square
		if rows*rows != len(in) {
			fmt.Printf("Warning: input matrix is not a perfect square, falling back to default\n")
			return def
		}

		// Initialize the matrix
		m := make([][]float64, rows)
		for i := range m {
			m[i] = make([]float64, rows)
		}

		// Fill the matrix with values from the input slice
		for i := 0; i < rows; i++ {
			for j := 0; j < rows; j++ {
				m[i][j] = in[i*rows+j]
			}
		}

		return m
	}
	return def
}

func (s *ImageFilter) Apply(i *image.Image) *image.Image {
	if s.Type == "" {
		return i
	}
	m, fn := Filters.Get(strings.ToLower(s.Type))
	if m != nil && fn != nil {
		fn(s, i, m)
		return i
	}
	fmt.Printf("Error: unknown filter type: %s\n", s.Type)
	return i
}

func NewImageFilter(typ string, options map[string]any) *ImageFilter {
	return &ImageFilter{
		Type:    typ,
		Options: options,
	}
}
