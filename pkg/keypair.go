package pkg

import (
	"bytes"
	"os"

	"github.com/go-faster/errors"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
)

// A key pair for asymmetric encryption
type KeyPair struct {
	Public []byte
	Secret []byte
}

// Derive a key pair from a seed file, which can hold arbitrary secret data
func (keyPair *KeyPair) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read key file")
	}

	seed := argon2.IDKey(data, []byte("sillysecrets"), 1, 64*1024, 4, 32)

	publicP, secretP, err := box.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		return errors.Wrap(err, "failed to generate key")
	}

	keyPair.Public = publicP[:]
	keyPair.Secret = secretP[:]

	return nil
}
