package cmd

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	pathm "path"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:     "dump <folder>",
	Aliases: []string{"u"},

	Short: "Dump all accessible secrets into a folder",

	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]

		storage := pkg.Storage{}
		if err := storage.Load(jsonPath); err != nil {
			return errors.Wrap(err, "failed to load storage")
		}

		keys := pkg.Keys{}
		if err := keys.Load(keyPaths); err != nil {
			return errors.Wrap(err, "failed to load keys")
		}

		for _, entry := range storage {
			msg, _, err := entry.Decrypt(keys)
			if err != nil {
				// this entry didn't belong to us, but that is not fatal
				slog.Debug("failed to decrypt",
					slog.String("entry", entry.Name),
					slog.String("error", err.Error()))
				continue
			}

			slog.Info("dumping", slog.String("entry", entry.Name))

			path := pathm.Join(dir, entry.Name)
			if err := os.MkdirAll(pathm.Dir(path), 0755); err != nil {
				return errors.Wrap(err, "failed to create dump directory")
			}

			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0400)
			if err != nil {
				return errors.Wrap(err, "failed to create dump file")
			}
			defer file.Close()

			if _, err := io.Copy(file, bytes.NewReader(msg)); err != nil {
				return errors.Wrapf(err, "failed to dump `%v`", entry.Name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
