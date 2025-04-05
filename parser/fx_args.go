package parser

import (
	"strconv"
	"strings"

	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/core/fx"
)

func parseNamedArgs(argsStr string, filter fx.Function, vars map[string]string) {
	// named args always show up in pairs (key = value) and there can be several on a single line
	// so we first have to identify the pairs, the left hand side is easy because it can't be a string
	// so there are no spaces to consider. once we get to the right hand side we do have to take
	// strings into account.
	lhs := ""
	for strings.Contains(argsStr, constants.ASSIGN) {
		// as long as there is a constants.STR_ASSIGN, there is yet another pair to resolve
		lhs = strings.TrimSpace(argsStr[:strings.Index(argsStr, constants.ASSIGN)])
		argsStr = strings.TrimSpace(argsStr[strings.Index(argsStr, constants.ASSIGN)+1:])
		inQuote := false
		idx := 0
		// now we have to find the first whitespace that isn't part of a string
		for _, c := range argsStr {
			complete := false
			switch c {
			case constants.CHAR_QUOTE:
				if inQuote && argsStr[idx-1] == constants.CHAR_ESCAPE {
					idx++
					continue // this is an escaped quote
				}
				inQuote = !inQuote // toggle quote status
			case constants.CHAR_SPACE, constants.CHAR_TAB:
				if !inQuote {
					val := argsStr[:idx]
					// found the right hand side end
					arg := fx.FunctionArg{
						Name:  lhs,
						Value: parseArgsValue(val, vars),
					}
					funcArgs := filter.GetArgs()
					for i, a := range funcArgs {
						if a.Name == lhs {
							funcArgs[i] = arg
							break
						}
					}
					argsStr = strings.TrimSpace(argsStr[idx+1:])
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
	if len(argsStr) > 0 {
		// there should be one final argument left
		arg := fx.FunctionArg{
			Name:  lhs,
			Value: parseArgsValue(argsStr, vars),
		}
		funcArgs := filter.GetArgs()
		for i, a := range funcArgs {
			if a.Name == lhs {
				funcArgs[i] = arg
				break
			}
		}
	}
}

func parseUnnamedArgs(argsStr string, filter fx.Function, vars map[string]string) error {
	argsList := filter.GetArgs()
	values := make([]any, len(argsList))
	valuesRaw := make([]any, len(argsList))
	inQuote := false
	inArg := false
	idx := 0
	argIdx := 0
	for i, c := range argsStr {
		switch c {
		case constants.CHAR_QUOTE:
			if inQuote && argsStr[i-1] == constants.CHAR_ESCAPE {
				continue // this is an escaped quote
			}
			inQuote = !inQuote
			if !inQuote && inArg {
				// we just finished a string
				val := strings.TrimSpace(argsStr[idx : i+1])
				values[argIdx] = parseArgsValue(val, vars)
				valuesRaw[argIdx] = val
				argIdx++
				idx = i + 1
				inArg = false
			}
		case constants.CHAR_SPACE, constants.CHAR_TAB:
			if !inQuote && inArg {
				val := strings.TrimSpace(argsStr[idx:i])
				values[argIdx] = parseArgsValue(val, vars)
				valuesRaw[argIdx] = val
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
	if argIdx < len(argsList) {
		val := strings.TrimSpace(argsStr[idx:])
		values[argIdx] = parseArgsValue(val, vars)
		valuesRaw[argIdx] = val
	}

	for i, arg := range argsList {
		arg.Value = values[i]
		argsList[i] = arg
	}
	return nil
}

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
	return strings.ReplaceAll(strings.Trim(value, constants.QUOTE), constants.ESCAPE+constants.QUOTE, constants.QUOTE)
}
