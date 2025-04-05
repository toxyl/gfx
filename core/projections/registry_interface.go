package projections

// ProjectionRegistry defines the interface for projection registration
type ProjectionRegistry interface {
	Register(p *RegistryProjection)
	Get(name string) *RegistryProjection
	List() []*RegistryProjection
}

// defaultRegistry is the default implementation of ProjectionRegistry
type defaultRegistry struct {
	projections map[string]*RegistryProjection
}

// NewRegistry creates a new projection registry
func NewRegistry() ProjectionRegistry {
	return &defaultRegistry{
		projections: make(map[string]*RegistryProjection),
	}
}

// Register registers a new projection
func (r *defaultRegistry) Register(p *RegistryProjection) {
	r.projections[p.Name()] = p
}

// Get returns a projection by name
func (r *defaultRegistry) Get(name string) *RegistryProjection {
	return r.projections[name]
}

// List returns all registered projections
func (r *defaultRegistry) List() []*RegistryProjection {
	result := make([]*RegistryProjection, 0, len(r.projections))
	for _, p := range r.projections {
		result = append(result, p)
	}
	return result
}

// Default is the default projection registry
var Default = NewRegistry()
