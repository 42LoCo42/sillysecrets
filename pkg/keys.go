package pkg

import (
	"log/slog"
	"os"
	pathm "path"

	"github.com/go-faster/errors"
)

// A map from base64-encoded public keys -> corresponding key pairs
type Keys map[string]KeyPair

// Load keys from a list of files and directories.
// The latter will NOT be traversed recursively.
func (keys *Keys) Load(paths []string) error {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			slog.Debug("failed to stat key", slog.String("path", path))
			continue
		}

		if info.IsDir() {
			slog.Debug("loading keys from", slog.String("directory", path))

			entries, err := os.ReadDir(path)
			if err != nil {
				slog.Debug("failed to read key", slog.String("directory", path))
				continue
			}

			for _, entry := range entries {
				file := pathm.Join(path, entry.Name())
				slog.Debug("loading key from", slog.String("file", file))
				if err := keys.loadFile(file); err != nil {
					slog.Debug("failed to read key", slog.String("file", file))
					continue
				}
			}
		} else {
			slog.Debug("loading key from", slog.String("file", path))
			if err := keys.loadFile(path); err != nil {
				slog.Debug("failed to read key", slog.String("file", path))
				continue
			}
		}
	}

	return nil
}

// Load a single key file into this map
func (keys *Keys) loadFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "failed to stat key file")
	}

	if info.IsDir() {
		slog.Debug("ignoring directory", slog.String("path", path))
		return nil
	}

	keyPair := KeyPair{}
	if err := keyPair.Load(path); err != nil {
		return errors.Wrap(err, "failed to load key from file")
	}

	(*keys)[Encode(keyPair.Public)] = keyPair
	return nil
}
