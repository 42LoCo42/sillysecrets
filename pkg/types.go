package sillysecrets

import (
	set "github.com/deckarep/golang-set/v2"
	"gopkg.in/yaml.v3"
)

type Quoted string

// MarshalYAML implements yaml.Marshaler.
func (s Quoted) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.DoubleQuotedStyle,
		Value: string(s),
	}, nil
}

func Unquote(qs []Quoted) []string {
	us := make([]string, len(qs))
	for i, q := range qs {
		us[i] = string(q)
	}
	return us
}

type Group struct {
	Key Quoted `yaml:"key,omitempty"`

	ContainsRaw []Quoted        `yaml:"contains,omitempty"`
	Contains    set.Set[string] `yaml:"-"`

	GrantsRaw []Quoted        `yaml:"grants,omitempty"`
	Grants    set.Set[string] `yaml:"-"`

	Secrets map[string]Quoted `yaml:"secrets,omitempty"`
}

type Groups map[string]Group
