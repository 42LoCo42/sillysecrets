package sillysecrets

import set "github.com/deckarep/golang-set/v2"

type Group struct {
	Key string `yaml:"key,omitempty"`

	ContainsRaw []string        `yaml:"contains,omitempty"`
	Contains    set.Set[string] `yaml:"-"`

	GrantsRaw []string        `yaml:"grants,omitempty"`
	Grants    set.Set[string] `yaml:"-"`

	Secrets map[string]string `yaml:"secrets,omitempty"`
}

type Groups map[string]Group
