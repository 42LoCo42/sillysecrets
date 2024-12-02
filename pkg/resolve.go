package sillysecrets

import "github.com/go-faster/errors"

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

// func ResolveToContains(groups Groups) (Groups, error) {
// 	return Resolve(groups, "grants",
// 		func(g Group) []string {
// 			return g.Grants.ToSlice()
// 		},
// 		func(n string, tg Group) {
// 			tg.Contains.Add(n)
// 		})
// }

func ResolveToGrants(groups Groups) (Groups, error) {
	return Resolve(groups, "contains",
		func(g Group) []string {
			return g.Contains.ToSlice()
		},
		func(n string, tg Group) {
			tg.Grants.Add(n)
		})
}
