package parser

import "strings"

func trimCommentsAndWhitespace(line string) string {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, STR_COMMENT) {
		return line // nothing to strip, so we return the line as is
	}
	if line[0] == CHAR_COMMENT {
		return "" // entire line is commented, so we strip everything
	}
	if !strings.Contains(line, STR_QUOTE) {
		return strings.TrimSpace(line[:strings.Index(line, STR_COMMENT)]) // there is no string in this line so we can simply strip the comment
	}

	// if we get here there is a string in the line and a # which could be part of the string, or the string is part of the comment
	// we have to scan each char until we find the first # that is not part of a string
	inQuote := false
	for i, c := range line {
		switch c {
		case CHAR_QUOTE:
			if inQuote && line[i-1] == CHAR_ESCAPE {
				continue // this is an escaped backtick
			}
			inQuote = !inQuote // toggle inQuote status
		case CHAR_COMMENT:
			if !inQuote {
				return strings.TrimSpace(line[:i]) // we're not in a string, so this is the comment position
			}
		}
	}

	// we shouldn't get here, let's error
	panic("stripping comments failed on: \n" + line + "\n")
}
