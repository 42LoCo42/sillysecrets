package sillysecrets_test

import (
	"log"
	"testing"

	"github.com/42LoCo42/sillysecrets"
	set "github.com/deckarep/golang-set/v2"
)

func Test(t *testing.T) {
	groups, err := sillysecrets.Load("example/sesi.yaml")
	if err != nil {
		t.Fatal(err)
	}

	log.Print("ResolveToContains")
	for n, g := range groups {
		have := g.Contains

		want := map[string]set.Set[string]{
			"admin:alice": set.NewSet(
				"machine:lazuli",
			),
			"machine:lazuli": set.NewSet(
				"user:bob",
			),
			"user:bob": set.NewSet[string](),
			"user:friend": set.NewSet(
				"user:bob",
			),
		}[n]

		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}

	log.Print("ResolveToGrants")
	for n, g := range groups {
		have := g.Grants

		want := map[string]set.Set[string]{
			"admin:alice": set.NewSet[string](),
			"machine:lazuli": set.NewSet(
				"admin:alice",
			),
			"user:bob": set.NewSet(
				"machine:lazuli",
				"user:friend",
			),
			"user:friend": set.NewSet[string](),
		}[n]

		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}

	log.Print("CollectSecrets")
	for n := range groups {
		have, err := sillysecrets.CollectSecrets(n, groups)
		if err != nil {
			t.Fatal(err)
		}

		want := map[string]set.Set[string]{
			"admin:alice": set.NewSet(
				"machine:lazuli.example",
				"user:bob.password",
			),
			"machine:lazuli": set.NewSet(
				"machine:lazuli.example",
				"user:bob.password",
			),
			"user:bob": set.NewSet(
				"user:bob.password",
			),
			"user:friend": set.NewSet(
				"user:bob.password",
			),
		}[n]

		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}

	log.Print("CollectKeys")
	for n := range groups {
		have, err := sillysecrets.CollectKeys(n, groups)
		if err != nil {
			t.Fatal(err)
		}

		want := map[string]set.Set[string]{
			"admin:alice": set.NewSet(
				groups["admin:alice"].Key,
			),
			"machine:lazuli": set.NewSet(
				groups["admin:alice"].Key,
				groups["machine:lazuli"].Key,
			),
			"user:bob": set.NewSet(
				groups["admin:alice"].Key,
				groups["machine:lazuli"].Key,
				groups["user:friend"].Key,
			),
			"user:friend": set.NewSet(
				groups["user:friend"].Key,
			),
		}[n]

		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}
}
