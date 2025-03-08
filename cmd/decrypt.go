package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:     "decrypt <secret>",
	Aliases: []string{"d"},

	Short: "Decrypt a secret to stdout",

	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: pkg.CobraValidSecrets,

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		storage := pkg.Storage{}
		if err := storage.Load(jsonPath); err != nil {
			return errors.Wrap(err, "failed to load storage")
		}

		entry, ok := storage[name]
		if !ok {
			return errors.Errorf("secret `%v` isn't created yet!", name)
		}

		keys := pkg.Keys{}
		if err := keys.Load(keyPaths); err != nil {
			return errors.Wrap(err, "failed to load keys")
		}

		dec, _, err := entry.Decrypt(keys)
		if err != nil {
			return errors.Wrap(err, "failed to decrypt entry")
		}

		if _, err := io.Copy(os.Stdout, bytes.NewReader(dec)); err != nil {
			return errors.Wrap(err, "failed to output entry")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
