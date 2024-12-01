package sillysecrets

import (
	"fmt"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
)

func Collect(
	name string,
	groups Groups,
	get func(g Group) []string,
	add func(n string, g Group, r set.Set[string]),
) (set.Set[string], error) {
	visited := set.NewSet[string]()
	results := set.NewSet[string]()

	var helper func(name string) error
	helper = func(name string) error {
		if visited.ContainsOne(name) {
			return nil
		}

		group, ok := groups[name]
		if !ok {
			return errors.Errorf("group %v not found", name)
		}

		visited.Add(name)

		for _, n := range get(group) {
			if err := helper(n); err != nil {
				return errors.Wrap(err, "could not collect subresults")
			}
		}

		add(name, group, results)
		results.Remove("")
		return nil
	}

	if err := helper(name); err != nil {
		return nil, err
	}
	return results, nil
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
			r.Add(string(g.Key))
		})
}
