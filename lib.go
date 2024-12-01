package sillysecrets

import (
	"fmt"
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

////////////////////////////////////////////////////////////////////////////////

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

	groups, err = ResolveToContains(groups)
	if err != nil {
		return nil, err
	}

	groups, err = ResolveToGrants(groups)
	if err != nil {
		return nil, err
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

////////////////////////////////////////////////////////////////////////////////

func Resolve(
	groups Groups,
	msg string,
	get func(g Group) []string,
	add func(n string, tg Group),
) (Groups, error) {
	for n, g := range groups {
		for _, tn := range get(g) {
			tg, ok := groups[tn]
			if !ok {
				return nil, errors.Errorf(
					"%v: %v: invalid target group %v",
					n, msg, tn)
			}

			add(n, tg)
			groups[tn] = tg
		}
	}

	return groups, nil
}

func ResolveToContains(groups Groups) (Groups, error) {
	return Resolve(groups, "grants",
		func(g Group) []string {
			return g.Grants.ToSlice()
		},
		func(n string, tg Group) {
			tg.Contains.Add(n)
		})
}

func ResolveToGrants(groups Groups) (Groups, error) {
	return Resolve(groups, "contains",
		func(g Group) []string {
			return g.Contains.ToSlice()
		},
		func(n string, tg Group) {
			tg.Grants.Add(n)
		})
}

////////////////////////////////////////////////////////////////////////////////

func Collect(
	name string,
	groups Groups,
	get func(g Group) []string,
	add func(n string, g Group, r set.Set[string]),
) (set.Set[string], error) {
	visited := set.NewSet[string]()

	var helper func(name string) (set.Set[string], error)
	helper = func(name string) (set.Set[string], error) {
		if visited.ContainsOne(name) {
			return nil, nil
		}

		group, ok := groups[name]
		if !ok {
			return nil, errors.Errorf("group %v not found", name)
		}

		visited.Add(name)
		results := set.NewSet[string]()

		for _, n := range get(group) {
			subresults, err := helper(n)
			if err != nil {
				return nil, err
			}

			results = results.Union(subresults)
		}

		add(name, group, results)
		results.Remove("")
		return results, nil
	}
	return helper(name)
}

func CollectSecrets(name string, groups Groups) (set.Set[string], error) {
	return Collect(name, groups,
		func(g Group) []string {
			return g.Contains.ToSlice()
		},
		func(n string, g Group, r set.Set[string]) {
			for secret := range g.Secrets {
				r.Add(fmt.Sprintf("%v.%v", n, secret))
			}
		})
}

func CollectKeys(name string, groups Groups) (set.Set[string], error) {
	return Collect(name, groups,
		func(g Group) []string {
			return g.Grants.ToSlice()
		},
		func(n string, g Group, r set.Set[string]) {
			r.Add(g.Key)
		})
}
