package sillysecrets

import (
	"os"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
)

const NOWRAP = "# -*- mode: yaml; truncate-lines: t; -*- vi: nowrap\n\n"

func Load(path string) (groups Groups, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	if err := yaml.NewDecoder(file).Decode(&groups); err != nil {
		return nil, errors.Wrap(err, "could not decode file")
	}

	for n, g := range groups {
		g.Contains = set.NewSet(Unquote(g.ContainsRaw)...)
		g.Grants = set.NewSet(Unquote(g.GrantsRaw)...)
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

	if _, err := file.WriteString(NOWRAP); err != nil {
		return errors.Wrap(err, "could not write nowrap magic to file")
	}

	if err := yaml.NewEncoder(file).Encode(groups); err != nil {
		return errors.Wrap(err, "could not encode file")
	}

	return nil
}
