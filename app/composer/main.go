package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/toxyl/gfx/parser"
)

func main() {
	var (
		fileIn      = flag.String("in", "", "(required) composition file")
		fileOut     = flag.String("out", "", "(required) output file (png or jpg)")
		fileOutGFXS = flag.String("gfxs", "", "(optional) path where to save parsed composition file (gfxs)")
		fileOutYAML = flag.String("yaml", "", "(optional) path where to save parsed composition file (yaml)")
	)

	flag.Parse()

	if strings.TrimSpace(*fileIn) == "" {
		fmt.Printf("no input file given!\n")
		flag.Usage()
		return
	}

	if strings.TrimSpace(*fileOut) == "" {
		fmt.Printf("no output file given!\n")
		flag.Usage()
		return
	}

	comp := parser.NewComposition("", 0, 0).LoadGFXS(*fileIn)
	f := *fileOut
	fl := strings.ToLower(f)

	if strings.HasSuffix(fl, ".png") {
		comp.Render().SaveAsPNG(f)
	}
	if strings.HasSuffix(fl, ".jpg") || strings.HasSuffix(fl, ".jpeg") {
		comp.Render().SaveAsJPG(f)
	}
	if strings.TrimSpace(*fileOutGFXS) != "" {
		comp.SaveGFXS(*fileOutGFXS)
	}
	if strings.TrimSpace(*fileOutYAML) != "" {
		comp.SaveYAML(*fileOutYAML)
	}
}
