package cmd

import (
	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var rekeyCmd = &cobra.Command{
	Use:     "rekey <secrets...>",
	Aliases: []string{"r"},

	Short: "Regenerate the internal shared key of some secrets",

	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: pkg.CobraValidSecrets,

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

		for _, name := range args {
			node, err := tree.Get(name)
			if err != nil {
				return errors.Wrapf(err, "failed to load node `%v`", name)
			}

			entry, ok := storage[name]
			if !ok {
				return errors.Errorf("secret `%v` isn't created yet!", name)
			}

			msg, _, err := entry.Decrypt(keys)
			if err != nil {
				return errors.Wrap(err, "failed to decrypt secret")
			}

			// create a completely new entry; we want a new shared key
			entry = &pkg.Entry{Name: name}
			if err := entry.Create(msg, node.AllKeys()); err != nil {
				return errors.Wrap(err, "failed to encrypt secret")
			}

			storage[name] = entry
		}

		if err := storage.Save(jsonPath); err != nil {
			return errors.Wrap(err, "failed to save storage")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rekeyCmd)
}
