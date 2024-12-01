package sillysecrets_test

import (
	"crypto/rand"
	"log"
	"path"
	"reflect"
	"testing"

	"github.com/42LoCo42/sillysecrets/pkg"
	set "github.com/deckarep/golang-set/v2"
)

const DIR = "../example"

func Test(t *testing.T) {
	groups, err := sillysecrets.Load(path.Join(DIR, "sesi.yaml"))
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
				string(groups["admin:alice"].Key),
			),
			"machine:lazuli": set.NewSet(
				string(groups["admin:alice"].Key),
				string(groups["machine:lazuli"].Key),
			),
			"user:bob": set.NewSet(
				string(groups["admin:alice"].Key),
				string(groups["machine:lazuli"].Key),
				string(groups["user:friend"].Key),
			),
			"user:friend": set.NewSet(
				string(groups["user:friend"].Key),
			),
		}[n]

		log.Printf("  %v: have %v, want %v", n, have, want)
		if !have.Equal(want) {
			t.Fatal(have, want)
		}
	}

	log.Print("Crypto")

	identities := sillysecrets.LoadIdentities([]string{DIR})
	log.Print("  Identities found:")
	for _, i := range identities {
		log.Print("    ", i)
	}

	for n := range groups {
		raw := make([]byte, 128)
		if _, err := rand.Read(raw); err != nil {
			t.Fatal(err)
		}

		enc, err := sillysecrets.Encrypt(raw, n, groups)
		if err != nil {
			t.Fatal(err)
		}

		dec, err := sillysecrets.Decrypt(enc, identities)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(raw, dec) {
			t.Fatal(enc, dec)
		}
	}
}
