package fx

// RegistryInterface defines the interface for the function registry
type RegistryInterface interface {
	// Register registers a new function with the registry
	Register(f Function) error

	// Get returns a function by name
	Get(name string) Function

	// List returns a list of all registered functions
	List() []string

	// GetMeta returns the metadata for a function
	GetMeta(name string) *FunctionMeta
}
