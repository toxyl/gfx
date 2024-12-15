package main

import (
	"fmt"
	"os"

	"github.com/toxyl/gfx/parser"
)

func main() {
	vsixPath, err := parser.GenerateVSIX()
	if err != nil {
		fmt.Println("Error generating VSCode extension:", err)
		os.Exit(1)
	}
	fmt.Print(vsixPath)
}
