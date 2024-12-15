package parser

import "strings"

func parseVarsSection(line string, vars map[string]string) {
	parts := strings.SplitN(line, STR_ASSIGN, 2)
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	vars[key] = value
}
