package sillysecrets_test

import (
	"log"
	"testing"

	"github.com/42LoCo42/sillysecrets"
)

func Must[X any](x X, err error) X {
	if err != nil {
		log.Panic(err)
	}

	return x
}

func Test(t *testing.T) {
	groups, err := sillysecrets.Load("example/sesi.yaml")
	if err != nil {
		t.Fatal(err)
	}

	log.Print("ResolveToContains")
	contains_want := Must(sillysecrets.Load("example/sesi.contains.yaml"))
	contains_have, err := sillysecrets.ResolveToContains(groups)
	if err != nil {
		t.Fatal(t)
	}
	for n := range contains_have {
		have := contains_have[n].Contains
		want := contains_want[n].Contains
		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}

	log.Print("ResolveToGrants")
	grants_want := Must(sillysecrets.Load("example/sesi.grants.yaml"))
	grants_have, err := sillysecrets.ResolveToGrants(groups)
	if err != nil {
		t.Fatal(t)
	}
	for n := range grants_have {
		have := grants_have[n].Grants
		want := grants_want[n].Grants
		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}
}
