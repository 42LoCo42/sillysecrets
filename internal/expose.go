package internal

import "filippo.io/age"

func ParseIdentitiesFile(name string) ([]age.Identity, error) {
	return parseIdentitiesFile(name)
}
