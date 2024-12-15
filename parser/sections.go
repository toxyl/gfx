package parser

import "strings"

func parseSection(line string) string {
	if strings.HasPrefix(line, STR_LBRACKET) && strings.HasSuffix(line, STR_RBRACKET) {
		return line[1 : len(line)-1]
	}
	return ""
}
