package filterchain

import (
	"strings"

	"github.com/toxyl/errors"
	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/filters"
	"github.com/toxyl/gfx/image"
)

type FilterChain []*filters.ImageFilter

func (fc *FilterChain) Save(file string) error {
	chain := []string{}
	for _, f := range *fc {
		chain = append(chain, f.String())
	}
	return flo.File(file).StoreString(strings.Join(chain, "\n"))
}

func (fc *FilterChain) Append(filter ...string) {
	for _, f := range filter {
		(*fc) = append(*fc, filters.Parse(f))
	}
}

func (fc *FilterChain) Apply(img *image.Image) *image.Image {
	for _, f := range *fc {
		img = f.Apply(img)
	}
	return img
}

func (fc *FilterChain) Load(file string) error {
	fChain := FilterChain{}
	f := flo.File(file)
	if !f.Exists() {
		return errors.Newf("file not found: %s", file)
	}
	for _, f := range strings.Split(f.AsString(), "\n") {
		if strings.TrimSpace(f) == "" {
			continue
		}
		fChain = append(fChain, filters.Parse(f))
	}
	*fc = fChain
	return nil
}

func New(filter ...string) *FilterChain {
	fc := &FilterChain{}
	fc.Append(filter...)
	return fc
}
