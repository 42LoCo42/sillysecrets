package cmd

import (
	"log"
	"strings"

	"filippo.io/age"
	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
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

var validSecretArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveError
	}

	return secrets().ToSlice(), cobra.ShellCompDirectiveNoFileComp
}

type LoadedSecret struct {
	RawName    string
	GroupName  string
	SecretName string
	Group      sillysecrets.Group
	Value      []byte
}

func loadSecret(args []string) LoadedSecret {
	s := LoadedSecret{}
	s.RawName = args[0]

	parts := strings.Split(s.RawName, ".")
	if len(parts) != 2 {
		log.Fatalf("invalid secret %v: must be in <group>.<secret> format", s.RawName)
	}

	s.GroupName = strings.TrimSpace(parts[0])
	s.SecretName = strings.TrimSpace(parts[1])
	if s.GroupName == "" || s.SecretName == "" {
		log.Fatalf("invalid secret %v: must be in <group>.<secret> format", s.RawName)
	}

	var ok bool
	s.Group, ok = groups()[s.GroupName]
	if !ok {
		log.Fatalf("group %v not found", s.GroupName)
	}

	if s.Group.Secrets == nil {
		s.Group.Secrets = map[string]sillysecrets.Quoted{}
	}

	enc, ok := s.Group.Secrets[s.SecretName]
	s.Value = []byte{}
	if ok {
		var err error
		s.Value, err = sillysecrets.Decrypt(enc, ids())
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not decrypt data"))
		}
	}

	return s
}
