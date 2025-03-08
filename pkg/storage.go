package pkg

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/go-faster/errors"
)

// A map of secret names -> entries holding their encrypted data and recipient information
type Storage map[string]*Entry

// Load the storage from a JSON file, optionally creating it
func (storage *Storage) Load(path string) error {
	slog.Debug("loading storage from", slog.String("file", path))

	if _, err := os.Stat(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "failed to stat storage file")
		}

		file, err := os.Create(path)
		if err != nil {
			return errors.Wrap(err, "failed to create storage file")
		}
		defer file.Close()

		if _, err := file.WriteString("{}"); err != nil {
			return errors.Wrap(err, "failed to initialize storage file")
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open storage file")
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(storage); err != nil {
		return errors.Wrap(err, "failed to load storage")
	}

	for name, entry := range *storage {
		entry.Name = name
	}

	return nil
}

// Save the storage to a JSON file, pruning entries with no recipients
func (storage *Storage) Save(path string) error {
	slog.Debug("saving storage to", slog.String("file", path))

	for name, entry := range *storage {
		if len(entry.Rcp) == 0 {
			slog.Debug("deleting empty entry", slog.String("entry", name))
			delete(*storage, name)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create storage file")
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	if err := enc.Encode(storage); err != nil {
		return errors.Wrap(err, "failed to save storage")
	}

	return nil
}
