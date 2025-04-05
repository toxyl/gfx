package fx

import (
	"image"

	"github.com/toxyl/gfx/core/color"
)

// ArgumentType represents the type of a function argument
type ArgumentType string

const (
	TypeInt    ArgumentType = "int"
	TypeFloat  ArgumentType = "float"
	TypeBool   ArgumentType = "bool"
	TypeString ArgumentType = "string"
	TypeColor  ArgumentType = "color"
)

// FunctionArg represents a function argument
type FunctionArg struct {
	Name        string
	Type        ArgumentType
	Value       any
	Default     any
	Required    bool
	Description string
	Min         float64
	Max         float64
	Step        float64
}

// FunctionMeta represents metadata about a function
type FunctionMeta struct {
	Name        string
	Description string
	Category    string
	Args        []FunctionArg
	Examples    []string
	Tags        []string
}

// FunctionRegistry is a map of function names to their implementations
type FunctionRegistry map[string]Effect

// ImageFunction represents a function that operates on an image
type ImageFunction interface {
	// ProcessPixel processes a single pixel
	ProcessPixel(x, y int, img *image.Image) (*color.Color64, error)
}
