package parser

import (
	"strings"

	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/core/fx"
)

type FX struct {
	Name  string
	Funcs []fx.Function
}

func (f *FX) String() string {
	conf := []string{}
	for _, filter := range f.Funcs {
		conf = append(conf, filter.Name())
	}
	return f.Name + constants.SPACE + constants.LBRACE + "\n" + constants.TAB + strings.Join(conf, "\n"+constants.TAB) + "\n" + constants.RBRACE + "\n"
}

func (f *FX) Append(filters ...fx.Function) *FX {
	f.Funcs = append(f.Funcs, filters...)
	return f
}

func (f *FX) Get() []fx.Function {
	return f.Funcs
}

func NewFX(name string) *FX {
	f := FX{
		Name:  name,
		Funcs: []fx.Function{},
	}
	return &f
}
