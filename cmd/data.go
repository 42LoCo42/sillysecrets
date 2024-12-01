package cmd

import (
	"filippo.io/age"
	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	set "github.com/deckarep/golang-set/v2"
)

var _groups sillysecrets.Groups

func groups() sillysecrets.Groups {
	if _groups == nil {
		_groups, _ = sillysecrets.Load(file)
	}
	return _groups
}

var _secrets set.Set[string]

func secrets() set.Set[string] {
	if _secrets == nil {
		_secrets = sillysecrets.AllSecrets(groups())
	}
	return _secrets
}

var _ids []age.Identity

func ids() []age.Identity {
	if _ids == nil {
		_ids = sillysecrets.LoadIdentities(idPaths)
	}
	return _ids
}
