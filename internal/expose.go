package internal

import "filippo.io/age"

func ParseIdentitiesFile(name string) ([]age.Identity, error) {
	return parseIdentitiesFile(name)
}

func ParseRecipient(arg string) (age.Recipient, error) {
	return parseRecipient(arg)
}
