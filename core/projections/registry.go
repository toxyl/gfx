package projections

import (
	"sync"
)

var (
	projections = make(map[string]*RegistryProjection)
	mu          sync.RWMutex
)

// Register registers a new projection
func Register(p *RegistryProjection) {
	mu.Lock()
	defer mu.Unlock()
	projections[p.Name()] = p
}

// Get returns a projection by name
func Get(name string) *RegistryProjection {
	mu.RLock()
	defer mu.RUnlock()
	return projections[name]
}

// List returns all registered projections
func List() []*RegistryProjection {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]*RegistryProjection, 0, len(projections))
	for _, p := range projections {
		result = append(result, p)
	}
	return result
}
