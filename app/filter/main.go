package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/toxyl/flo"
	"github.com/toxyl/gfx/filterchain"
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
		chain     filterChain
		fileIn    = flag.String("i", "", "input file")
		fileOut   = flag.String("o", "", "output file")
		fileChain = flag.String("c", "", "filter chain file (if present the filter chain will be loaded from this file instead of the -f flags, if not present the -f flags will be used to create the file)")
	)
	flag.Var(&chain, "f", "filters (use multiple -f flags for a filter chain)")
	flag.Parse()

	if fileIn != nil && *fileIn != "" && fileOut != nil && *fileOut != "" {
		fChain := filterchain.New()
		appendFilters := true
		saveChain := false
		if fileChain != nil && *fileChain != "" {
			f := flo.File(*fileChain)
			if f.Exists() {
				if err := fChain.Load(f.Path()); err != nil {
					fmt.Printf("Failed to load filter chain: %s\n", err.Error())
					return
				}
				appendFilters = false
				fmt.Printf("Filter chain loaded from file.\n")
			} else {
				saveChain = true
			}
		}
		if appendFilters {
			fChain.Append(chain...)
		}
		if saveChain {
			fChain.Save(*fileChain)
			fmt.Printf("Filter chain saved to file.\n")
		}
		fChain.Apply(image.NewFromFile(*fileIn)).SaveAsPNG(*fileOut)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "filters" {
		fmt.Printf("Available filters\n-----------------\n%s\n", strings.Join(filters.Examples, "\n"))
		return
	}

	flag.Usage()
}
