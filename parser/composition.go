package parser

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/math"
)

type Composition struct {
	Name   string          `yaml:"name,omitempty"`
	Width  int             `yaml:"width,omitempty"`
	Height int             `yaml:"height,omitempty"`
	Layers []*Layer        `yaml:"layers,omitempty"`
	Color  *hsla.HSLA      `yaml:"color,omitempty"`
	Crop   *Crop           `yaml:"crop,omitempty"`
	Resize *Resize         `yaml:"resize,omitempty"`
	Filter *CompiledFilter `yaml:"filter,omitempty"`
}

func (c *Composition) String() string {
	maxLenOp := math.MaxLenStr(COMPOSITION...)
	spf := fmt.Sprintf
	spfPad := func(l int, s string) string { return spf("%-*s", l, s) }
	color := STR_COMMENT + " no " + COMP_COLOR + " defined"
	filter := STR_COMMENT + " no " + COMP_FILTER + " defined"
	resize := STR_COMMENT + " no " + COMP_RESIZE + " defined"
	crop := STR_COMMENT + " no " + COMP_CROP + " defined"
	name := STR_COMMENT + " no " + COMP_NAME + " defined"
	width := STR_COMMENT + " no " + COMP_WIDTH + " defined"
	height := STR_COMMENT + " no " + COMP_HEIGHT + " defined"
	if c.Color != nil {
		color = spf("%s %s hsla%s%f %f %f %f%s", spfPad(maxLenOp, COMP_COLOR), STR_ASSIGN, STR_LPAREN, c.Color.H(), c.Color.S(), c.Color.L(), c.Color.A(), STR_RPAREN)
	}
	if c.Filter != nil {
		filter = spf("%s %s %s", spfPad(maxLenOp, COMP_FILTER), STR_ASSIGN, c.Filter.Name)
	}
	if c.Crop != nil {
		crop = spf("%s %s %d %d %d %d", spfPad(maxLenOp, COMP_CROP), STR_ASSIGN, c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H)
	}
	if c.Resize != nil {
		resize = spf("%s %s %d %d", spfPad(maxLenOp, COMP_RESIZE), STR_ASSIGN, c.Resize.W, c.Resize.H)
	}
	if c.Name != "" {
		name = spf("%s %s %s", spfPad(maxLenOp, COMP_NAME), STR_ASSIGN, STR_QUOTE+c.Name+STR_QUOTE)
	}
	if c.Width != 0 {
		width = spf("%s %s %d", spfPad(maxLenOp, COMP_WIDTH), STR_ASSIGN, c.Width)
	}
	if c.Height != 0 {
		height = spf("%s %s %d", spfPad(maxLenOp, COMP_HEIGHT), STR_ASSIGN, c.Height)
	}
	filters := []string{}
	layers := []string{}
	if c.Layers != nil {
		hasCrop, hasResize, hasOffset, hasFilter := false, false, false, false
		for _, l := range c.Layers {
			if l == nil {
				continue
			}
			if l.Crop != nil {
				hasCrop = true
			}
			if l.Resize != nil {
				hasResize = true
			}
			if l.Offset != nil {
				hasOffset = true
			}
			if l.Filter != nil {
				hasFilter = true
			}
		}

		for _, l := range c.Layers {
			if l == nil {
				continue
			}

			layers = append(layers, l.String(hasCrop, hasResize, hasOffset, hasFilter))
			if l.Filter != nil {
				filters = append(filters, l.Filter.String())
			}
		}
	}
	if c.Filter != nil {
		filters = append(filters, c.Filter.String())
	}
	return spf(
		`%s%s%s
# this is a compiled filter that doesn't contain any variables anymore

%s%s%s
%s

%s%s%s
%s 
%s
%s
%s 
%s
%s
%s

%s%s%s
%s
`,
		STR_LBRACKET, strings.ToUpper(SECTION_VARS), STR_RBRACKET,
		STR_LBRACKET, strings.ToUpper(SECTION_FILTERS), STR_RBRACKET,
		strings.Join(filters, "\n"),
		STR_LBRACKET, strings.ToUpper(SECTION_COMPOSITION), STR_RBRACKET,
		name,
		width,
		height,
		color,
		filter,
		crop,
		resize,
		STR_LBRACKET, strings.ToUpper(SECTION_LAYERS), STR_RBRACKET,
		strings.Join(layers, "\n"),
	)
}

func NewComposition(name string, w, h int) *Composition {
	c := Composition{
		Name:   name,
		Width:  w,
		Height: h,
		Layers: []*Layer{},
		Color:  nil,
		Crop:   nil,
		Resize: nil,
		Filter: nil,
	}
	return &c
}
