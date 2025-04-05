package fx

import (
	"image"
)

// Apply applies a function to an image
func Apply(img image.Image, f Effect) image.Image {
	if f == nil {
		return img
	}
	result, _ := f.Apply(img)
	return result
}

// Register registers a function with the default registry
func Register(f Effect) error {
	return DefaultRegistry.Register(f)
}

// Get returns a function by name from the default registry
func Get(name string) Effect {
	return DefaultRegistry.Get(name)
}

// List returns a list of all registered functions
func List() []string {
	return DefaultRegistry.List()
}

// GetMeta returns the metadata for a function
func GetMeta(name string) *FunctionMeta {
	return DefaultRegistry.GetMeta(name)
}
