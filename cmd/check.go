package cmd

import (
	"log/slog"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},

	Short: "Ensure congruency between structure and storage",
	Long: `Ensure congruency between structure and storage.
This compares the recipient set of every storage entry
with the expected key set as derived from the structure tree.
If there is a mismatch, the entry will be adjusted accordingly`,

	Args: cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		tree := pkg.Tree{}
		if err := tree.Load(yamlPath); err != nil {
			return errors.Wrap(err, "failed to load tree")
		}

		storage := pkg.Storage{}
		if err := storage.Load(jsonPath); err != nil {
			return errors.Wrap(err, "failed to load storage")
		}

		keys := pkg.Keys{}
		if err := keys.Load(keyPaths); err != nil {
			return errors.Wrap(err, "failed to load keys")
		}

		for name, entry := range storage {
			node, err := tree.Get(name)
			if err != nil {
				return errors.Wrapf(err, "failed to get node `%v`", name)
			}

			have := entry.Available()
			want := node.AllKeys()

			if have.Equal(want.Set) {
				continue
			}

			slog.Debug("fixing incongruency",
				slog.String("entry", name),
				slog.Any("have", have),
				slog.Any("want", want))

			_, shared, err := entry.Decrypt(keys)
			if err != nil {
				return errors.Wrapf(err, "failed to decrypt entry `%v`", name)
			}

			// only encrypt for new recipients; leave msg/enc unchanged
			if err := entry.EncryptRcp(shared, want); err != nil {
				return errors.Wrapf(err, "failed to encrypt entry `%v`", name)
			}
		}

		if err := storage.Save(jsonPath); err != nil {
			return errors.Wrap(err, "failed to save storage")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
