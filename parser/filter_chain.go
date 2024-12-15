package parser

import (
	"strconv"
	"strings"

	"github.com/toxyl/errors"
	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/image"
)

func parseChain(str string) *ImageFilter {
	args := strings.Split(str, STR_LPAREN)
	name := strings.TrimSpace(args[0])
	args = strings.Split(strings.TrimSuffix(args[1], STR_RPAREN), STR_SPACE)
	options := map[string]any{}
	for _, a := range args {
		e := strings.Split(a, STR_ASSIGN)
		if len(e) != 2 {
			continue
		}
		k := strings.TrimSpace(e[0])
		v := strings.TrimSpace(e[1])
		if name == convolution.Meta.Name && k == convolution.Meta.NameOf(0) {
			m := []float64{}
			for _, e := range strings.Split(v, STR_COMMA) {
				e = strings.TrimSpace(e)
				if f, err := strconv.ParseFloat(e, 64); err == nil {
					m = append(m, f)
				} else {
					m = append(m, 0)
				}
			}
			options[k] = m
		} else if f, err := strconv.ParseFloat(v, 64); err == nil {
			options[k] = f
		} else if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			options[k] = i
		} else {
			options[k] = v
		}
	}
	return NewImageFilter(name, options)
}

type FilterChain []*ImageFilter

func (fc *FilterChain) Load(file string) error {
	fChain := FilterChain{}
	f := flo.File(file)
	if !f.Exists() {
		return errors.Newf("file not found: %s", file)
	}
	for _, f := range strings.Split(f.AsString(), "\n") {
		f = strings.TrimSpace(f)
		if f == "" || f[0] == CHAR_COMMENT {
			continue
		}
		fChain = append(fChain, parseChain(f))
	}
	*fc = fChain
	return nil
}

func (fc *FilterChain) Save(file string) error {
	chain := []string{}
	for _, f := range *fc {
		chain = append(chain, f.String(true))
	}
	return flo.File(file).StoreString(strings.Join(chain, "\n"))
}

func (fc *FilterChain) Append(filter ...string) {
	for _, f := range filter {
		(*fc) = append(*fc, parseChain(f))
	}
}

func (fc *FilterChain) Apply(img *image.Image) *image.Image {
	out := img.Clone()
	for _, f := range *fc {
		out = f.Apply(out)
	}
	return out
}

func NewFilterChain(filter ...string) *FilterChain {
	fc := &FilterChain{}
	fc.Append(filter...)
	return fc
}
