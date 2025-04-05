package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [src] <dst>", filepath.Base(os.Args[0]))
		fmt.Printf("`src` will overwritten with the formatted result `dst` is not given.")
		return
	}
	f := os.Args[1]
	comp, err := parser.ParseComposition(fs.LoadString(f))
	if err != nil {
		panic("Failed to parse composition: " + err.Error())
	}

	if len(os.Args) == 3 {
		f = os.Args[2]
	}
	comp.Save(f)
}
