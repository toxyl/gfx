package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

type Crop struct {
	X int
	Y int
	W int
	H int
}

func (c *Crop) String() string {
	return fmt.Sprintf("%s = %4d %4d %4d %4d", constants.COMP_CROP, c.X, c.Y, c.W, c.H)
}

type Resize struct {
	W int
	H int
}

func (r *Resize) String() string {
	return fmt.Sprintf("%s = %4d %4d", constants.COMP_RESIZE, r.W, r.H)
}

func parseComments(line string) string {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, constants.COMMENT) {
		return line // nothing to strip, so we return the line as is
	}
	if line[0] == constants.CHAR_COMMENT {
		return "" // entire line is commented, so we strip everything
	}
	if !strings.Contains(line, constants.QUOTE) {
		return strings.TrimSpace(line[:strings.Index(line, constants.COMMENT)]) // there is no string in this line so we can simply strip the comment
	}

	// if we get here there is a string in the line and a # which could be part of the string, or the string is part of the comment
	// we have to scan each char until we find the first # that is not part of a string
	inQuote := false
	for i, c := range line {
		switch c {
		case constants.CHAR_QUOTE:
			if inQuote && line[i-1] == constants.CHAR_ESCAPE {
				continue // this is an escaped backtick
			}
			inQuote = !inQuote // toggle inQuote status
		case constants.CHAR_COMMENT:
			if !inQuote {
				return strings.TrimSpace(line[:i]) // we're not in a string, so this is the comment position
			}
		}
	}

	// we shouldn't get here, let's error
	panic("stripping comments failed on: \n" + line + "\n")
}

func parseColor(value string) *color.HSL {
	value = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(value)), "hsla"+constants.LPAREN), constants.RPAREN)
	parts := strings.Split(value, constants.SPACE)
	if len(parts) != 4 {
		panic("composition color must be given as `hsla" + constants.LPAREN + "hue" + constants.SPACE + "sat" + constants.SPACE + "lum" + constants.SPACE + "alpha" + constants.RPAREN + "`")
	}
	h, _ := strconv.ParseFloat(parts[0], 64)
	s, _ := strconv.ParseFloat(parts[1], 64)
	l, _ := strconv.ParseFloat(parts[2], 64)
	a, _ := strconv.ParseFloat(parts[3], 64)
	c, _ := color.NewHSL(h, s, l, a)
	return c
}

func parseCrop(value string) Crop {
	parts := strings.Split(value, constants.SPACE)
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	w, _ := strconv.Atoi(parts[2])
	h, _ := strconv.Atoi(parts[3])
	return Crop{X: x, Y: y, W: w, H: h}
}

func parseResize(value string) Resize {
	parts := strings.Split(value, constants.SPACE)
	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])
	return Resize{W: w, H: h}
}

func parseCompositionSection(line string, comp *Composition, filters map[string]*FX) {
	parts := strings.SplitN(line, constants.ASSIGN, 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	switch key {
	case constants.COMP_FILTER:
		comp.Filter.Append(filters[value].Funcs...)
	case constants.COMP_CROP:
		crop := parseCrop(value)
		comp.Crop = &crop
	case constants.COMP_RESIZE:
		resize := parseResize(value)
		comp.Resize = &resize
	case constants.COMP_COLOR:
		comp.Color = parseColor(value)
	case constants.COMP_WIDTH:
		comp.Width, _ = strconv.Atoi(value)
	case constants.COMP_HEIGHT:
		comp.Height, _ = strconv.Atoi(value)
	}
}

func parseVarsSection(line string, vars map[string]string) {
	parts := strings.SplitN(line, constants.ASSIGN, 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	vars[key] = value
}

func parseFilter(script []string, vars map[string]string, fns map[string]*FX) (fn []fx.Effect, needsHoisting bool) {
	parsedFilters := []fx.Effect{}
	lkw := len(constants.KEYWORD_USE) + 1
	for _, line := range script {
		line = strings.TrimSpace(line)
		filterType := line[:strings.Index(line, constants.LPAREN)]
		if strings.EqualFold(filterType, constants.KEYWORD_USE) {
			filterName := strings.TrimSpace(line[lkw : len(line)-1])
			if _, ok := fns[filterName]; !ok {
				return nil, true
			}
			parsedFilters = append(parsedFilters, fns[filterName].Get()...)
			continue
		}
		fnArgs := strings.TrimSpace(line[strings.Index(line, constants.LPAREN)+1 : strings.Index(line, constants.RPAREN)])

		filter := fx.NewBaseFunction(
			filterType,
			"",
			color.New(0, 0, 0, 1),
			[]fx.FunctionArg{},
		)

		if fnArgs != "" {
			if strings.Contains(fnArgs, constants.ASSIGN) {
				parseNamedArgs(fnArgs, filter, vars)
			} else {
				if err := parseUnnamedArgs(fnArgs, filter, vars); err != nil {
					fmt.Printf("Error parsing unnamed arguments for filter %s: %v\n", filterType, err)
				}
			}
		}

		parsedFilters = append(parsedFilters, filter)
	}
	return parsedFilters, false
}

func parseFilterBlock(line string, scanner *bufio.Scanner) (string, []string) {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, constants.LBRACE) {
		return "", nil // this is not a filter block
	}
	idxLBrace := strings.Index(line, constants.LBRACE)
	idxLast := len(line) - 1
	filterName := strings.TrimSpace(line[:idxLBrace])
	filterData := ""
	if line[idxLast] != constants.CHAR_RBRACE {
		// this is a multiline block, so we first have to find the end,
		// select everything up to there and replace linebreaks with spaces,
		// so we don't have to worry about the amount of filters defined per line,
		// in the next step we will split them into single filters
		for scanner.Scan() {
			nextLine := parseComments(scanner.Text())
			if nextLine == "" {
				continue
			}
			// the line might contain an STR_RBRACE in a string, so we have to check char by char
			inQuote := false
			isEnd := false
			for i, c := range nextLine {
				shouldBreak := false
				switch c {
				case constants.CHAR_QUOTE:
					if inQuote && nextLine[i-1] == constants.CHAR_ESCAPE {
						continue // this is an escaped quote
					}
				case constants.CHAR_RBRACE:
					if !inQuote {
						// this is the end of the filter block
						isEnd = true
						shouldBreak = true
					}
				}
				if shouldBreak {
					break
				}
			}
			if isEnd {
				break
			}
			filterData += constants.SPACE + nextLine
		}
	} else {
		// this is a single line filter block
		filterData = strings.TrimSpace(line[idxLBrace+1 : idxLast])
	}

	var filterLines []string
	// there might be multiple filters in the block
	// and the arguments might contain strings
	// which can contain constants.STR_LPAREN or constants.STR_RPAREN,
	// so we have to check char by char and
	// if we encounter a constants.STR_RPAREN that is not
	// part of a string, we have a full filter and
	// can append it to the list
	inQuote := false
	idxCurrent := 0
	for i, c := range filterData {
		switch c {
		case constants.CHAR_QUOTE:
			if inQuote && filterData[i-1] == constants.CHAR_ESCAPE {
				continue // this is an escaped quote char
			}
			inQuote = !inQuote // toggle quote status
		case constants.CHAR_RPAREN:
			if !inQuote {
				// we have a complete filter, add it
				filterLines = append(filterLines, strings.TrimSpace(filterData[idxCurrent:i+1]))
				idxCurrent = i + 1
			}
		}
	}
	return filterName, filterLines
}

func parseLayer(line string, filters map[string]*FX) Layer {
	parts := strings.Fields(strings.TrimSpace(line))
	alpha, _ := strconv.ParseFloat(parts[1], 64)
	src := strings.Join(parts[3:], constants.SPACE) // we join everything left so paths can contain spaces
	return *NewLayer(parts[0], alpha, filters[parts[2]], src)
}

func parseSection(line string) string {
	if strings.HasPrefix(line, constants.LBRACKET) && strings.HasSuffix(line, constants.RBRACKET) {
		return line[1 : len(line)-1]
	}
	return ""
}

func ParseComposition(content string) (*Composition, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	comp := Composition{
		Width:  0,
		Height: 0,
		Layers: []*Layer{},
		Crop:   &Crop{},
		Resize: &Resize{},
		Filter: NewFX("c" + fmt.Sprint(time.Now().UnixMilli())), // create a default filter for the composition
	}
	var currentSection string
	vars := make(map[string]string)
	filters := make(map[string]*FX)
	lines := struct {
		vars           []string
		filtersSection string
		filters        map[string][]string
		composition    []string
		layers         []string
	}{
		vars:           []string{},
		filtersSection: "",
		filters:        map[string][]string{},
		composition:    []string{},
		layers:         []string{},
	}

	// extract and sort lines
	for scanner.Scan() {
		line := parseComments(scanner.Text())
		if line == "" {
			continue
		}
		if s := parseSection(line); s != "" {
			currentSection = s
			continue
		}
		switch strings.ToUpper(currentSection) {
		case constants.SECTION_VARS:
			lines.vars = append(lines.vars, line)
		case constants.SECTION_FILTERS:
			lines.filtersSection += line + "\n"
		case constants.SECTION_COMPOSITION:
			lines.composition = append(lines.composition, line)
		case constants.SECTION_LAYERS:
			lines.layers = append(lines.layers, line)
		}
	}

	// process vars
	for _, line := range lines.vars {
		parseVarsSection(line, vars)
	}
	for k, v := range vars {
		comp.Vars = append(comp.Vars, &Var{
			Name:  k,
			Value: v,
		})
	}

	// process filter blocks
	scanner2 := bufio.NewScanner(strings.NewReader(lines.filtersSection))
	for scanner2.Scan() {
		name, script := parseFilterBlock(scanner2.Text(), scanner2)
		if name != "" {
			lines.filters[name] = script
		}
	}

	// process filters
	n := len(lines.filters)
	passesWithoutChange := 0
	parsed := make(map[string]bool, n)
	for n > 0 {
		if passesWithoutChange > 2 {
			// we've tried several times, but there's nothing changing, so we must be stuck in a circular reference
			unparsed := []string{}
			for k := range lines.filters {
				if _, ok := parsed[k]; !ok {
					unparsed = append(unparsed, k)
				}
			}
			fmt.Printf("ERROR: Your script got stuck in a loop!\nDo you have a circular reference with `use` calls? \nCheck these filters: \n- %s\n", strings.Join(unparsed, "\n- "))
			os.Exit(1)
		}
		changes := 0
		for name, script := range lines.filters {
			if parsed[name] {
				continue // already parsed it, we're good
			}
			fns, needsHoisting := parseFilter(script, vars, filters)
			if needsHoisting {
				// this one has a use statement for a yet undefined filter, try again later
				continue
			}
			parsed[name] = true
			filters[name] = NewFX(name).Append(fns...)
			n--
			changes++
		}
		if changes == 0 {
			passesWithoutChange++
		}
	}

	// process composition
	for _, line := range lines.composition {
		parseCompositionSection(line, &comp, filters)
	}

	// process layers
	for _, line := range lines.layers {
		layer := parseLayer(line, filters)
		comp.Layers = append(comp.Layers, &layer)
	}

	return &comp, nil
}
