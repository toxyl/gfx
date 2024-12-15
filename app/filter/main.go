package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/parser"
)

type filterChain []string

func (c *filterChain) String() string {
	return strings.Join(*c, ", ")
}

func (c *filterChain) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func MakeExample(m *parser.MetaData) string {
	sb := strings.Builder{}
	sb.WriteString(m.Name + parser.STR_LPAREN)
	l := len(m.Args)
	for i, a := range m.Args {
		sb.WriteString(
			fmt.Sprintf("%s%s%v", a.Name, parser.STR_ASSIGN, a.Default),
		)
		if i < l-1 {
			sb.WriteString(parser.STR_SPACE)
		}
	}
	sb.WriteString(parser.STR_RPAREN)
	return sb.String()
}

func main() {
	var (
		chain     filterChain
		showList  = flag.Bool("list", false, "if provided a list with examples of all available filters will be printed, all other flags will be ignored")
		fileIn    = flag.String("in", "", "input file")
		fileOut   = flag.String("out", "", "output file")
		fileChain = flag.String("chain", "", "filter chain file (if present the filter chain will be loaded from this file instead of the -f flags, if not present the -f flags will be used to create the file)")
	)
	flag.Var(&chain, "f", "filters (use multiple -f flags for a filter chain)")
	flag.Parse()

	if *showList {
		fmt.Printf("Available filters\n-----------------\n")
		for _, f := range *parser.Filters {
			fmt.Printf("%s\n", MakeExample(f.Meta))
		}
		return
	}

	if fileIn != nil && *fileIn != "" && fileOut != nil && *fileOut != "" {
		filterChain := parser.NewFilterChain()
		appendFilters := true
		saveChain := false
		if fileChain != nil && *fileChain != "" {
			f := flo.File(*fileChain)
			if f.Exists() {
				if len(chain) > 0 {
					fmt.Printf("Warning: you can't use filters defined with `-f` flags if you also provide an existing filter chain file using the `-c` flag, ignoring all `-f` flags.\n")
				}
				if err := filterChain.Load(f.Path()); err != nil {
					fmt.Printf("Failed to load filter chain: %s\n", err.Error())
					return
				}
				appendFilters = false
			} else {
				saveChain = true
			}
		}

		if appendFilters {
			filterChain.Append(chain...)
		}
		if saveChain {
			filterChain.Save(*fileChain)
			fmt.Printf("Filter chain saved to %s.\n", *fileOut)
		}
		res := filterChain.Apply(image.NewFromFile(*fileIn))
		ft := strings.ToLower(*fileOut)
		if strings.HasSuffix(ft, ".png") {
			res.SaveAsPNG(*fileOut)
			return
		}
		if strings.HasSuffix(ft, ".jpg") || strings.HasSuffix(ft, ".jpeg") {
			res.SaveAsJPG(*fileOut)
			return
		}
		return
	}

	flag.Usage()
}
