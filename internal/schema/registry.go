package schema

import (
	"fmt"
	"sync"
)

// Registry holds a collection of named schemas and resolves them by name.
type Registry struct {
	mu      sync.RWMutex
	schemas map[string]*Schema
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{schemas: make(map[string]*Schema)}
}

// Register adds a schema to the registry.
// Returns an error if a schema with the same name already exists.
func (r *Registry) Register(s *Schema) error {
	if s == nil {
		return fmt.Errorf("cannot register nil schema")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schemas[s.Name]; exists {
		return fmt.Errorf("schema %q already registered", s.Name)
	}
	r.schemas[s.Name] = s
	return nil
}

// Get retrieves a schema by name. Returns nil, false if not found.
func (r *Registry) Get(name string) (*Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.schemas[name]
	return s, ok
}

// Names returns all registered schema names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.schemas))
	for n := range r.schemas {
		names = append(names, n)
	}
	return names
}

// MatchFirst tries each registered schema against raw and returns the first
// schema name and field map that matches, or ("", nil) if none match.
func (r *Registry) MatchFirst(raw string) (string, map[string]string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name, s := range r.schemas {
		if fields := s.Match(raw); fields != nil {
			return name, fields
		}
	}
	return "", nil
}
