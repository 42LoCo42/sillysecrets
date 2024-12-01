package sillysecrets

import set "github.com/deckarep/golang-set/v2"

type Group struct {
	Key string

	ContainsRaw []string        `yaml:"contains"`
	Contains    set.Set[string] `yaml:"-"`

	GrantsRaw []string        `yaml:"grants"`
	Grants    set.Set[string] `yaml:"-"`

	Secrets map[string]string
}

type Groups map[string]Group
