package parser

import (
	"strings"
)

func parseNamedArgs(args string, filter *ImageFilter, vars map[string]string) {
	// named args always show up in pairs (key = value) and there can be several on a single line
	// so we first have to identify the pairs, the left hand side is easy because it can't be a string
	// so there are no spaces to consider. once we get to the right hand side we do have to take
	// strings into account.
	lhs := ""
	for strings.Contains(args, STR_ASSIGN) {
		// as long as there is a STR_ASSIGN, there is yet another pair to resolve
		lhs = strings.TrimSpace(args[:strings.Index(args, STR_ASSIGN)])
		args = strings.TrimSpace(args[strings.Index(args, STR_ASSIGN)+1:])
		inQuote := false
		idx := 0
		// now we have to find the first whitespace that isn't part of a string
		for _, c := range args {
			complete := false
			switch c {
			case CHAR_QUOTE:
				if inQuote && args[idx-1] == CHAR_ESCAPE {
					idx++
					continue // this is an escaped quote
				}
				inQuote = !inQuote // toggle quote status
			case CHAR_SPACE, CHAR_TAB:
				if !inQuote {
					val := args[:idx]
					// found the right hand side end
					filter.Options[lhs] = parseArgsValue(val, vars) // set the option
					args = strings.TrimSpace(args[idx+1:])
					idx = 0
					complete = true
				}
			}
			if complete {
				break
			}
			idx++
		}
	}
	if len(args) > 0 {
		// there should be one final argument left
		filter.Options[lhs] = parseArgsValue(args, vars) // set the option
	}
}
