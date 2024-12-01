package sillysecrets

import (
	"os"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
)

type Group struct {
	Key string

	ContainsRaw []string        `yaml:"contains"`
	Contains    set.Set[string] `yaml:"-"`

	GrantsRaw []string        `yaml:"grants"`
	Grants    set.Set[string] `yaml:"-"`

	Secrets map[string]string
}

type Groups map[string]Group

func Load(path string) (groups Groups, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	if err := yaml.NewDecoder(file).Decode(&groups); err != nil {
		return nil, errors.Wrap(err, "could not decode file")
	}

	for n, g := range groups {
		g.Contains = set.NewSet(g.ContainsRaw...)
		g.Grants = set.NewSet(g.GrantsRaw...)
		groups[n] = g
	}

	return groups, nil
}

func Save(path string, groups Groups) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}

	if err := yaml.NewEncoder(file).Encode(groups); err != nil {
		return errors.Wrap(err, "could not encode file")
	}

	return nil
}

func ResolveToContains(groups Groups) (Groups, error) {
	for n, g := range groups {
		for _, tn := range g.Grants.ToSlice() {
			tg, ok := groups[tn]
			if !ok {
				return nil, errors.Errorf(
					"%v: grants: invalid target group %v",
					n, tn)
			}

			tg.Contains.Add(n)
			groups[tn] = tg
		}
	}

	return groups, nil
}

func ResolveToGrants(groups Groups) (Groups, error) {
	for n, g := range groups {
		for _, tn := range g.Contains.ToSlice() {
			tg, ok := groups[tn]
			if !ok {
				return nil, errors.Errorf(
					"%v: contains: invalid target group %v",
					n, tn)
			}

			tg.Grants.Add(n)
			groups[tn] = tg
		}
	}

	return groups, nil
}
