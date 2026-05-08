// Package plugin provides a lightweight extension point that allows
// external checkers to participate in drift detection.
package plugin

import (
	"fmt"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Checker is the interface that every plugin must implement.
type Checker interface {
	// Name returns a short, unique identifier for the plugin.
	Name() string

	// Check inspects the supplied drift results and may append extra
	// results or annotate existing ones. It returns the (possibly
	// modified) slice or an error.
	Check(results []drift.Result) ([]drift.Result, error)
}

// Registry holds registered Checker implementations.
type Registry struct {
	checkers []Checker
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a Checker to the registry. It returns an error if a
// checker with the same name has already been registered.
func (r *Registry) Register(c Checker) error {
	for _, existing := range r.checkers {
		if existing.Name() == c.Name() {
			return fmt.Errorf("plugin %q is already registered", c.Name())
		}
	}
	r.checkers = append(r.checkers, c)
	return nil
}

// Run executes every registered Checker in registration order, passing
// the accumulated results from one checker into the next.
func (r *Registry) Run(results []drift.Result) ([]drift.Result, error) {
	for _, c := range r.checkers {
		var err error
		results, err = c.Check(results)
		if err != nil {
			return results, fmt.Errorf("plugin %q: %w", c.Name(), err)
		}
	}
	return results, nil
}

// Names returns the names of all registered checkers in order.
func (r *Registry) Names() []string {
	names := make([]string, len(r.checkers))
	for i, c := range r.checkers {
		names[i] = c.Name()
	}
	return names
}
