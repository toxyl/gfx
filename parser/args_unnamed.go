package parser

import (
	"strings"
)

func parseUnnamedArgs(args string, filter *ImageFilter, vars map[string]string) error {
	m, _ := Filters.Get(filter.Type)
	keys := m.ArgNames()
	values := make([]any, len(keys))
	inQuote := false
	inArg := false
	idx := 0
	argIdx := 0
	for i, c := range args {
		switch c {
		case CHAR_QUOTE:
			if inQuote && args[i-1] == CHAR_ESCAPE {
				continue // this is an escaped quote
			}
			inQuote = !inQuote
			if !inQuote && inArg {
				// we just finished a string
				val := strings.TrimSpace(args[idx : i+1])
				values[argIdx] = parseArgsValue(val, vars)
				argIdx++
				idx = i + 1
				inArg = false
			}
		case CHAR_SPACE, CHAR_TAB:
			if !inQuote && inArg {
				val := strings.TrimSpace(args[idx:i])
				values[argIdx] = parseArgsValue(val, vars)
				argIdx++
				idx = i
				inArg = false
			}
		default:
			if !inArg {
				inArg = true
			}
		}
	}
	if argIdx < len(keys) {
		val := strings.TrimSpace(args[idx:])
		values[argIdx] = parseArgsValue(val, vars)
	}

	for i := 0; i < len(keys); i++ {
		filter.Options[keys[i]] = values[i]
	}
	return nil
}
