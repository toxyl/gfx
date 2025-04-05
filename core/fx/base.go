package fx

import (
	"errors"
	"image"

	"github.com/toxyl/gfx/core/color"
)

// Error types
var (
	ErrInvalidImage    = errors.New("invalid image")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrOutOfBounds     = errors.New("coordinates out of bounds")
)

// Effect represents a base interface for all image processing functions
type Effect interface {
	// Name returns the name of the function
	Name() string

	// Description returns a description of what the function does
	Description() string

	// Apply applies the function to the given image
	Apply(img image.Image) (image.Image, error)

	// GetColorModel returns the color model used by this function
	GetColorModel() *color.Color64

	// ValidateArgs validates the function arguments
	ValidateArgs(args ...any) error

	// GetArgs returns the function arguments
	GetArgs() []FunctionArg
}

// BaseFunction provides common functionality for all functions
type BaseFunction struct {
	name        string
	description string
	colorModel  *color.Color64
	args        []FunctionArg
}

// NewBaseFunction creates a new base function
func NewBaseFunction(name, description string, colorModel *color.Color64, args []FunctionArg) *BaseFunction {
	return &BaseFunction{
		name:        name,
		description: description,
		colorModel:  colorModel,
		args:        args,
	}
}

// Name implements the Function interface
func (f *BaseFunction) Name() string {
	return f.name
}

// Description implements the Function interface
func (f *BaseFunction) Description() string {
	return f.description
}

// GetColorModel implements the Function interface
func (f *BaseFunction) GetColorModel() *color.Color64 {
	return f.colorModel
}

// GetArgs implements the Function interface
func (f *BaseFunction) GetArgs() []FunctionArg {
	return f.args
}

// ValidateArgs implements the Function interface
func (f *BaseFunction) ValidateArgs(args ...any) error {
	if len(args) != len(f.args) {
		return ErrInvalidArgument
	}
	return nil
}

// Apply implements the Function interface
func (f *BaseFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, ErrInvalidImage
	}
	return img, nil
}
