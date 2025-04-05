package fx

import (
	"fmt"
	"sync"
)

// Registry is the global function registry
type Registry struct {
	mu        sync.RWMutex
	functions FunctionRegistry
}

// NewRegistry creates a new function registry
func NewRegistry() *Registry {
	return &Registry{
		functions: make(FunctionRegistry),
	}
}

// Register implements RegistryInterface
func (r *Registry) Register(f Effect) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := f.Name()
	if _, exists := r.functions[name]; exists {
		return fmt.Errorf("function %s already registered", name)
	}

	r.functions[name] = f
	return nil
}

// Get implements RegistryInterface
func (r *Registry) Get(name string) Effect {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.functions[name]
}

// List implements RegistryInterface
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// GetMeta implements RegistryInterface
func (r *Registry) GetMeta(name string) *FunctionMeta {
	f := r.Get(name)
	if f == nil {
		return nil
	}

	return &FunctionMeta{
		Name:        name,
		Description: f.Description(),
		Args:        f.GetArgs(),
	}
}

// DefaultRegistry is the default global registry
var DefaultRegistry = NewRegistry()
