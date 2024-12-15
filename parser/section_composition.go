package parser

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/image"
)

func parseCompositionSection(line string, comp *Composition, filters map[string]*CompiledFilter) {
	parts := strings.SplitN(line, STR_ASSIGN, 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	switch key {
	case COMP_FILTER:
		comp.Filter.Append(filters[value].Filters...)
	case COMP_CROP:
		crop := parseCrop(value)
		comp.Crop = &crop
	case COMP_RESIZE:
		resize := parseResize(value)
		comp.Resize = &resize
	case COMP_COLOR:
		comp.Color = parseColor(value)
	case COMP_NAME:
		comp.Name = strings.Trim(value, STR_QUOTE)
	case COMP_WIDTH:
		comp.Width, _ = strconv.Atoi(value)
	case COMP_HEIGHT:
		comp.Height, _ = strconv.Atoi(value)
	}
}

func ParseComposition(content string) (*Composition, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	comp := Composition{
		Name:   "",
		Width:  0,
		Height: 0,
		Layers: []*Layer{},
		Crop:   &Crop{},
		Resize: &Resize{},
		Filter: NewCompiledFilter("compFilter"),
	}
	var currentSection string
	vars := make(map[string]string)
	fltrs := make(map[string]*CompiledFilter)

	for scanner.Scan() {
		line := trimCommentsAndWhitespace(scanner.Text())
		if line == "" {
			continue
		}
		if s := parseSection(line); s != "" {
			currentSection = s
			continue
		}
		switch strings.ToUpper(currentSection) {
		case SECTION_VARS:
			parseVarsSection(line, vars)
		case SECTION_FILTERS:
			filterName, filterLines := parseFilterBlock(line, scanner)
			if filterName != "" {
				fltrs[filterName] = NewCompiledFilter(filterName).Append(parseFilters(filterLines, vars, fltrs)...)
			}
		case SECTION_COMPOSITION:
			parseCompositionSection(line, &comp, fltrs)
		case SECTION_LAYERS:
			layer := parseLayer(line, fltrs)
			comp.Layers = append(comp.Layers, &layer)
		}
	}

	return &comp, nil
}

func (c *Composition) LoadYAML(path string) *Composition {
	if err := flo.File(path).LoadYAML(&c); err != nil {
		panic("failed to load composition from YAML: " + err.Error())
	}
	return c
}

func (c *Composition) Load(path string) *Composition {
	str := ""
	if err := flo.File(path).LoadString(&str); err != nil {
		panic("failed to load composition: " + err.Error())
	}
	comp, err := ParseComposition(str)
	if err != nil {
		panic("failed to parse composition: " + err.Error())
	}
	c = comp
	return c
}

func (c *Composition) SaveYAML(path string) *Composition {
	if err := flo.File(path).StoreYAML(&c); err != nil {
		panic("failed to save composition as YAML: " + err.Error())
	}
	return c
}

func (c *Composition) SaveGFXS(path string) *Composition {
	if err := flo.File(path).StoreString(c.String()); err != nil {
		panic("failed to save composition: " + err.Error())
	}
	return c
}

func (c *Composition) Render() *image.Image {
	w, h := c.Width, c.Height
	res := image.New(w, h)
	if c.Color != nil {
		res.FillHSLA(0, 0, w, h, c.Color)
	}
	numLayers := len(c.Layers)
	for i := numLayers - 1; i >= 0; i-- {
		l := c.Layers[i]
		scaled := l.Render(w, h)
		res.Draw(
			scaled,
			0, 0, w, h,
			0, 0, w, h,
			blend.BlendMode(l.BlendMode),
			l.Alpha,
		)
	}
	if c.Filter != nil {
		for _, filter := range c.Filter.Get() {
			if filter != nil {
				res = filter.Apply(res)
			}
		}
	}
	if c.Crop != nil && c.Crop.W > 0 && c.Crop.H > 0 {
		res = res.Crop(c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H, true)
	}
	if c.Resize != nil && c.Resize.W > 0 && c.Resize.H > 0 {
		res = res.Resize(c.Resize.W, c.Resize.H)
	}
	return res.Resize(w, h)
}
