package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [GFXScript file] [output image file] <[image file 1] .. [image file N]>", filepath.Base(os.Args[0]))
		return
	}
	fileIn := os.Args[1]
	fileOut := os.Args[2]
	images := []string{}

	if len(os.Args) > 3 {
		images = os.Args[3:]
	}

	if strings.TrimSpace(fileIn) == "" {
		fmt.Printf("no input file given!\n")
		flag.Usage()
		return
	}

	if strings.TrimSpace(fileOut) == "" {
		fmt.Printf("no output file given!\n")
		flag.Usage()
		return
	}

	t := time.Now()
	comp, err := parser.ParseComposition(fs.LoadString(fileIn))
	if err != nil {
		panic("Failed to parse composition: " + err.Error())
	}
	f := fileOut
	fl := strings.ToLower(f)
	comp.Save("test-parsed.gfxs")
	fmt.Printf("Composition parsed in %s.\n", time.Since(t).String())
	t = time.Now()
	if strings.HasSuffix(fl, ".png") {
		img, err := comp.Render(images...)
		if err == nil {
			err = img.SavePNG(f)
			if err != nil {
				panic("Failed to save PNG: " + err.Error())
			}
		}
	}
	if strings.HasSuffix(fl, ".jpg") || strings.HasSuffix(fl, ".jpeg") {
		img, err := comp.Render(images...)
		if err == nil {
			err = img.SaveJPEG(f, 90)
			if err != nil {
				panic("Failed to save JPEG: " + err.Error())
			}
		}
	}
	fmt.Printf("Rendered and saved in %s.\n", time.Since(t).String())
}
