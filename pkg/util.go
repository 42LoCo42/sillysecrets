package sillysecrets

import (
	"fmt"

	set "github.com/deckarep/golang-set/v2"
)

func AllSecrets(groups Groups) set.Set[string] {
	secrets := set.NewSet[string]()

	for n, g := range groups {
		for s := range g.Secrets {
			secrets.Add(fmt.Sprintf("%v.%v", n, s))
		}
	}

	return secrets
}
