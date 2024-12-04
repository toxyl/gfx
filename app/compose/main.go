package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/gfx/composition"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:   %s [composition file] [output path]\n", filepath.Base(os.Args[0]))
		fmt.Printf("Example: %s composition.yaml output.png\n", filepath.Base(os.Args[0]))
		return
	}
	composition.
		Load(os.Args[1]).     // load the composition from a file,
		Render().             // render it and
		SaveAsPNG(os.Args[2]) // save as PNG
}
