package filters

import (
	"strconv"
	"strings"
)

func Parse(str string) *ImageFilter {
	args := strings.Split(str, "::")
	name := strings.TrimSpace(args[0])
	options := map[string]any{}
	for _, a := range args[1:] {
		e := strings.Split(a, "=")
		if len(e) != 2 {
			continue
		}
		k := strings.TrimSpace(e[0])
		v := strings.TrimSpace(e[1])
		if k == optMtr {
			m := []float64{}
			for _, e := range strings.Split(v, ",") {
				e = strings.TrimSpace(e)
				if f, err := strconv.ParseFloat(e, 64); err == nil {
					m = append(m, f)
				} else {
					m = append(m, 0)
				}
			}
			options[k] = m
		} else if f, err := strconv.ParseFloat(v, 64); err == nil {
			options[k] = f
		} else if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			options[k] = i
		} else {
			options[k] = v
		}
	}
	return NewImageFilter(name, options)
}
