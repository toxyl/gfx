package parser

import (
	"bufio"
	"fmt"
	"strings"
)

func parseFilterBlock(line string, scanner *bufio.Scanner) (string, []string) {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, STR_LBRACE) {
		return "", nil // this is not a filter block
	}
	idxLBrace := strings.Index(line, STR_LBRACE)
	idxLast := len(line) - 1
	filterName := strings.TrimSpace(line[:idxLBrace])
	filterData := ""
	if line[idxLast] != CHAR_RBRACE {
		// this is a multiline block, so we first have to find the end,
		// select everything up to there and replace linebreaks with spaces,
		// so we don't have to worry about the amount of filters defined per line,
		// in the next step we will split them into single filters
		for scanner.Scan() {
			nextLine := trimCommentsAndWhitespace(scanner.Text())
			if nextLine == "" {
				continue
			}
			// the line might contain an STR_RBRACE in a string, so we have to check char by char
			inQuote := false
			isEnd := false
			for i, c := range nextLine {
				switch c {
				case CHAR_QUOTE:
					if inQuote && nextLine[i-1] == CHAR_ESCAPE {
						continue // this is an escaped quote
					}
				case CHAR_RBRACE:
					if !inQuote {
						// this is the end of the filter block
						isEnd = true
						break
					}
				}
			}
			if isEnd {
				break
			}
			filterData += STR_SPACE + nextLine
		}
	} else {
		// this is a single line filter block
		filterData = strings.TrimSpace(line[idxLBrace+1 : idxLast])
	}

	var filterLines []string
	// there might be multiple filters in the block
	// and the arguments might contain strings
	// which can contain STR_LPAREN or STR_RPAREN,
	// so we have to check char by char and
	// if we encounter a STR_RPAREN that is not
	// part of a string, we have a full filter and
	// can append it to the list
	inQuote := false
	idxCurrent := 0
	for i, c := range filterData {
		switch c {
		case CHAR_QUOTE:
			if inQuote && filterData[i-1] == CHAR_ESCAPE {
				continue // this is an escaped quote char
			}
			inQuote = !inQuote // toggle quote status
		case CHAR_RPAREN:
			if !inQuote {
				// we have a complete filter, add it
				filterLines = append(filterLines, strings.TrimSpace(filterData[idxCurrent:i+1]))
				idxCurrent = i + 1
			}
		}
	}
	return filterName, filterLines
}

func parseFilters(lines []string, vars map[string]string, fltrs map[string]*CompiledFilter) []*ImageFilter {
	parsedFilters := []*ImageFilter{}
	lkw := len(KEYWORD_USE) + 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		filterType := line[:strings.Index(line, STR_LPAREN)]
		if strings.EqualFold(filterType, KEYWORD_USE) {
			filterName := strings.TrimSpace(line[lkw : len(line)-1])
			parsedFilters = append(parsedFilters, fltrs[filterName].Get()...)
			continue
		}
		filterArgs := strings.TrimSpace(line[strings.Index(line, STR_LPAREN)+1 : strings.Index(line, STR_RPAREN)])

		filter := &ImageFilter{
			Type:    filterType,
			Options: make(map[string]any),
		}

		if filterArgs != "" {
			if strings.Contains(filterArgs, STR_ASSIGN) {
				parseNamedArgs(filterArgs, filter, vars)
			} else {
				if err := parseUnnamedArgs(filterArgs, filter, vars); err != nil {
					fmt.Printf("Error parsing unnamed arguments for filter %s: %v\n", filterType, err)
				}
			}
		}

		parsedFilters = append(parsedFilters, filter)
	}
	return parsedFilters
}
