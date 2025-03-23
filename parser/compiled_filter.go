package parser

import (
	"strings"
)

type CompiledFilter struct {
	Name    string         `yaml:"name,omitempty"`
	Filters []*ImageFilter `yaml:"filters,omitempty"`
}

func (f *CompiledFilter) String() string {
	conf := []string{}
	for _, filter := range f.Filters {
		conf = append(conf, filter.String(false))
	}
	return f.Name + STR_SPACE + STR_LBRACE + STR_SPACE + strings.Join(conf, STR_SPACE) + STR_SPACE + STR_RBRACE
}

func (f *CompiledFilter) Append(filters ...*ImageFilter) *CompiledFilter {
	f.Filters = append(f.Filters, filters...)
	return f
}

func (f *CompiledFilter) Get() []*ImageFilter {
	return f.Filters
}

func NewCompiledFilter(name string) *CompiledFilter {
	f := CompiledFilter{
		Name:    name,
		Filters: []*ImageFilter{},
	}
	return &f
}
