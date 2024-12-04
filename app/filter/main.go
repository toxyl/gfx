package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/toxyl/gfx/filters"
	"github.com/toxyl/gfx/image"
)

type filterChain []string

func (c *filterChain) String() string {
	return strings.Join(*c, ", ")
}

func (c *filterChain) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func main() {
	var (
		chain   filterChain
		fileIn  = flag.String("i", "", "input file")
		fileOut = flag.String("o", "", "output file")
	)
	flag.Var(&chain, "f", "filters (use multiple -f flags for a filter chain)")
	flag.Parse()

	if fileIn != nil && *fileIn != "" && fileOut != nil && *fileOut != "" {
		fChain := []*filters.ImageFilter{}
		for _, f := range chain {
			args := strings.Split(f, "::")
			name := args[0]
			options := map[string]any{}
			for _, a := range args[1:] {
				e := strings.Split(a, "=")
				if len(e) != 2 {
					continue
				}
				k := e[0]
				v := e[1]
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					options[k] = f
				} else if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					options[k] = i
				} else {
					options[k] = v
				}
			}
			fChain = append(fChain, filters.NewImageFilter(name, options))
		}
		img := image.NewFromFile(*fileIn)
		for _, f := range fChain {
			img = f.Apply(img)
		}
		img.SaveAsPNG(*fileOut)
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "filters" {
		fmt.Printf("Available filters\n-----------------\n%s\n", strings.Join(filters.EXAMPLES, "\n"))
		return
	}
	flag.Usage()
}
