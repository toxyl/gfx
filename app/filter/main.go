package main

import (
	"flag"
	"fmt"
	"os"
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
			fChain = append(fChain, filters.Parse(f))
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
