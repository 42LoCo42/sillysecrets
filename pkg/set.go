package pkg

import (
	"sort"

	"github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
)

// A set of strings
type Set struct {
	mapset.Set[string]
}

// Convert this set to a sorted YAML list
func (s Set) MarshalYAML() (any, error) {
	result := s.ToSlice()
	sort.Strings(result)
	return result, nil
}

// Fill this set with values from a YAML list
func (s *Set) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return errors.New("set must be defined as YAML list")
	}

	s.Set = mapset.NewSet[string]()

	for _, x := range value.Content {
		if x.Kind != yaml.ScalarNode {
			return errors.New("set item must be a string")
		}

		s.Add(x.Value)
	}

	return nil
}
