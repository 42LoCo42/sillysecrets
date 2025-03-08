package pkg

import (
	"crypto/rand"

	"github.com/go-faster/errors"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

///// secret key cryptography /////

// Generate a shared key for symmetric encryption
func GenKey() (key []byte, err error) {
	key = make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, errors.Wrap(err, "failed to generate key")
	}

	return key, nil
}

// Encrypt some data with a shared key
func Encrypt(msg []byte, key []byte) (enc string, err error) {
	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		return "", errors.Wrap(err, "failed to generate nonce")
	}

	// prepend nonce to encrypted data
	return Encode(secretbox.Seal(
		nonce, msg,
		(*[24]byte)(nonce),
		(*[32]byte)(key))), nil
}

// Decrypt some data with a shared key
func Decrypt(enc string, key []byte) (msg []byte, err error) {
	raw, err := Decode(enc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode data")
	}

	// data begins after 24 bytes of nonce
	msg, ok := secretbox.Open(
		nil, raw[24:],
		(*[24]byte)(raw),
		(*[32]byte)(key))
	if !ok {
		return nil, errors.Wrap(err, "failed to decrypt data")
	}

	return msg, nil
}

///// public key cryptography /////

// Encrypt some data for a public key
func EncryptTo(msg []byte, public []byte) (enc string, err error) {
	raw, err := box.SealAnonymous(nil, msg, (*[32]byte)(public), rand.Reader)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt data")
	}

	return Encode(raw), nil
}

// Decrypt some data with a key pair
func DecryptAs(enc string, keyPair KeyPair) (msg []byte, err error) {
	raw, err := Decode(enc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode data")
	}

	msg, ok := box.OpenAnonymous(
		nil, raw,
		(*[32]byte)(keyPair.Public),
		(*[32]byte)(keyPair.Secret))
	if !ok {
		return nil, errors.New("failed to decrypt data")
	}

	return msg, nil
}
