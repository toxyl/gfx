package parser

import (
	"strconv"
	"strings"
)

func parseArgsValue(value string, vars map[string]string) any {
	if val, ok := vars[value]; ok {
		return parseArgsValue(val, vars)
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return float64(i)
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return strings.ReplaceAll(strings.Trim(value, STR_QUOTE), STR_ESCAPE+STR_QUOTE, STR_QUOTE)
}
