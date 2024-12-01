package sillysecrets

import (
	"bytes"
	"io"
	"log"
	"os"
	pathm "path"

	"filippo.io/age"
	"github.com/42LoCo42/sillysecrets/internal"
	"github.com/42LoCo42/z85m"
	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
)

func Encrypt(raw []byte, name string, groups Groups) (string, error) {
	keys, err := CollectKeys(name, groups)
	if err != nil {
		return "", errors.Wrap(err, "could not collect keys")
	}

	recp, err := ParseRecipients(keys)
	if err != nil {
		return "", errors.Wrap(err, "could not parse keys")
	}

	encb := bytes.Buffer{}
	encw, err := age.Encrypt(&encb, recp...)
	if err != nil {
		return "", errors.Wrap(err, "could not prepare encryption")
	}

	if _, err := encw.Write(raw); err != nil {
		return "", errors.Wrap(err, "could not encrypt data")
	}

	encw.Close()

	enc, err := z85m.Encode(encb.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "could not encode data")
	}

	return string(enc), nil
}

func Decrypt(enc string, identities []age.Identity) ([]byte, error) {
	encb, err := z85m.Decode([]byte(enc))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode data")
	}

	encr, err := age.Decrypt(bytes.NewReader(encb), identities...)
	if err != nil {
		return nil, errors.Wrap(err, "could not decrypt data")
	}

	dec, err := io.ReadAll(encr)
	if err != nil {
		return nil, errors.Wrap(err, "could not read decrypted data")
	}

	return dec, nil
}

func ParseRecipients(keys set.Set[string]) ([]age.Recipient, error) {
	recp := []age.Recipient{}

	for _, k := range keys.ToSlice() {
		r, err := internal.ParseRecipient(k)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse recipient `%v`", k)
		}

		recp = append(recp, r)
	}

	return recp, nil
}

func LoadIdentities(idPaths []string) []age.Identity {
	ids := []age.Identity{}

	var helper func(path string)
	helper = func(path string) {
		if err := func() error {
			info, err := os.Stat(path)
			if err != nil {
				return errors.Wrap(err, "could not get file information")
			}

			if info.IsDir() {
				entries, err := os.ReadDir(path)
				if err != nil {
					return errors.Wrap(err, "could not read directory")
				}

				for _, entry := range entries {
					helper(pathm.Join(path, entry.Name()))
				}
			} else {
				file, err := os.Open(path)
				if err != nil {
					return errors.Wrap(err, "could not open file")
				}
				defer file.Close()

				subids, err := internal.ParseIdentitiesFile(path)
				if err != nil {
					return errors.Wrapf(err, "could not parse identities in %v", path)
				}

				ids = append(ids, subids...)
			}

			return nil
		}(); err != nil {
			log.Printf("WARN: %v", err)
		}
	}

	for _, path := range idPaths {
		helper(path)
	}

	return ids
}
