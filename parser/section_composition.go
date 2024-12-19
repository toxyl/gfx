package parser

import (
	"bufio"
	"strconv"
	"strings"
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
