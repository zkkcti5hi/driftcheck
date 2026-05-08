// Package ignore provides support for loading and evaluating drift ignore rules.
// Rules can suppress specific drift fields for named services, reducing noise
// from intentional or expected differences.
package ignore

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Rule describes a single ignore entry: a service name and the list of drift
// fields (e.g. "image", "env", "ports") that should be suppressed for it.
type Rule struct {
	Service string   `yaml:"service"`
	Fields  []string `yaml:"fields"`
}

// Set is a collection of ignore rules loaded from a file.
type Set struct {
	Rules []Rule `yaml:"ignore"`
}

// Load reads an ignore-rules file from path. If path is empty or the file does
// not exist, an empty Set is returned without error so callers can treat a
// missing file as "no ignores".
func Load(path string) (*Set, error) {
	if path == "" {
		return &Set{}, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Set{}, nil
	}
	if err != nil {
		return nil, err
	}

	var s Set
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Suppressed reports whether the given field should be ignored for the named
// service. The match is case-sensitive on both service name and field name.
func (s *Set) Suppressed(service, field string) bool {
	for _, r := range s.Rules {
		if r.Service != service {
			continue
		}
		for _, f := range r.Fields {
			if f == field {
				return true
			}
		}
	}
	return false
}
