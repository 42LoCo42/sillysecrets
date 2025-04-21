package cmd

import (
	"log/slog"
	"strings"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var delMissing bool

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
				if delMissing {
					slog.Warn("deleting storage entry with missing node", slog.String("entry", name))
					delete(storage, name)
				} else {
					slog.Warn("storage entry is missing node", slog.String("entry", name))
				}

				continue
			}

			have := entry.Available()
			want := node.AllKeys()

			if have.Equal(want.Set) {
				continue
			}

			args := []any{slog.String("entry", name)}

			if add := want.Difference(have.Set); !add.IsEmpty() {
				args = append(args, slog.String("add",
					wrapColor("32", strings.Join(add.ToSlice(), " "))))
			}

			if del := have.Difference(want.Set); !del.IsEmpty() {
				args = append(args, slog.String("del",
					wrapColor("31", strings.Join(del.ToSlice(), " "))))
			}

			slog.Warn("fixing incongruency", args...)

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

	checkCmd.Flags().BoolVarP(&delMissing,
		"delete", "D",
		false,
		"Delete storage entries with no corresponding tree nodes")
}
