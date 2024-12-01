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
