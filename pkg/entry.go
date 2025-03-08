package pkg

import (
	"log/slog"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
)

// The representation of an encrypted secret that is part of the storage
type Entry struct {
	Name string `json:"-"` // Entry name, corresponding to its place in the Storage map

	Enc string            `json:"enc"` // encrypted data
	Rcp map[string]string `json:"rcp"` // map of recipient public key -> shared key encrypted to that public key
}

// Return all available recipients
func (entry *Entry) Available() Set {
	return Set{mapset.NewSetFromMapKeys(entry.Rcp)}
}

// Create a fresh Entry by encrypting a message for some recipients
func (entry *Entry) Create(msg []byte, recipients Set) error {
	shared, err := GenKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate shared key")
	}

	return entry.Encrypt(msg, shared, recipients)
}

// Encrypt a message with the specified shared key and store the result in this entry
func (entry *Entry) EncryptMsg(msg []byte, shared []byte) (err error) {
	entry.Enc, err = Encrypt(msg, shared)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt data")
	}

	return nil
}

// Ensure that this entry holds the shared key for the set of recipients.
// Note that this does NOT change its encrypted data; use EncryptMsg or Encrypt for that!
func (entry *Entry) EncryptRcp(shared []byte, recipients Set) (err error) {
	if entry.Rcp == nil {
		entry.Rcp = map[string]string{}
	}

	have := entry.Available().Set
	want := recipients.Set

	for _, publicB64 := range have.Difference(want).ToSlice() {
		slog.Debug("deleting key",
			slog.String("entry", entry.Name),
			slog.String("key", publicB64))

		delete(entry.Rcp, publicB64)
	}

	for _, publicB64 := range want.Difference(have).ToSlice() {
		slog.Debug("adding key",
			slog.String("entry", entry.Name),
			slog.String("key", publicB64))

		public, err := Decode(publicB64)
		if err != nil {
			return errors.Wrapf(err, "failed to decode public key `%v`", publicB64)
		}

		sharedEnc, err := EncryptTo(shared, public)
		if err != nil {
			return errors.Wrapf(err, "failed to encrypt to public key `%v`", publicB64)
		}

		entry.Rcp[publicB64] = sharedEnc
	}

	return nil
}

// Encrypt a message with a given shared key & for a set of recipients.
// This is a combined version of EncryptMsg & EncryptRcp
func (entry *Entry) Encrypt(msg []byte, shared []byte, recipients Set) (err error) {
	if err := entry.EncryptMsg(msg, shared); err != nil {
		return err
	}

	if err := entry.EncryptRcp(shared, recipients); err != nil {
		return err
	}

	return nil
}

// Decrypt the entry using a set of possible keys.
// If there is no macthing recipient for any of the keys, decryption will fail.
func (entry *Entry) Decrypt(keys Keys) (msg []byte, shared []byte, err error) {
	for publicB64, sharedEnc := range entry.Rcp {
		keyPair, ok := keys[publicB64]
		if !ok {
			continue
		}

		slog.Debug("decrypting",
			slog.String("entry", entry.Name),
			slog.String("key", publicB64))

		shared, err := DecryptAs(sharedEnc, keyPair)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "`%v`: failed to decrypt shared key", entry.Name)
		}

		msg, err := Decrypt(entry.Enc, shared)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "`%v`: failed to decrypt message", entry.Name)
		}

		return msg, shared, nil
	}

	return nil, nil, errors.Errorf("`%v`: no matching key", entry.Name)
}
