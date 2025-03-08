package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit <secret>",
	Aliases: []string{"e"},

	Short: "Edit or create a secret",
	Long: `Edit or create a secret.
If stdin is a pipe, it will be read into the secret, overwriting it.
Otherwise, vipe(1) will be started and given the current value of the secret
(or nothing if it was just created) as input.`,

	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: pkg.CobraValidSecrets,

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tree := pkg.Tree{}
		if err := tree.Load(yamlPath); err != nil {
			return errors.Wrap(err, "failed to load tree")
		}

		node, err := tree.Get(name)
		if err != nil {
			return errors.Wrapf(err, "failed to load node `%v`", name)
		}

		storage := pkg.Storage{}
		if err := storage.Load(jsonPath); err != nil {
			return errors.Wrap(err, "failed to load storage")
		}

		entry, ok := storage[name]
		if !ok {
			entry = new(pkg.Entry)
			storage[name] = entry

			msg, err := editor([]byte{})
			if err != nil {
				return err
			}

			if err := entry.Create(msg, node.AllKeys()); err != nil {
				return errors.Wrap(err, "failed to create entry")
			}
		} else {
			keys := pkg.Keys{}
			if err := keys.Load(keyPaths); err != nil {
				return errors.Wrap(err, "failed to load keys")
			}

			msg, shared, err := entry.Decrypt(keys)
			if err != nil {
				return errors.Wrap(err, "failed to decrypt entry")
			}

			new, err := editor(msg)
			if err != nil {
				return err
			}

			if !bytes.Equal(msg, new) {
				if err := entry.Encrypt(new, shared, node.AllKeys()); err != nil {
					return errors.Wrap(err, "failed to encrypt entry")
				}
			}
		}

		if err := storage.Save(jsonPath); err != nil {
			return errors.Wrap(err, "failed to save storage")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editor(old []byte) (new []byte, err error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat stdin")
	}

	buf := bytes.Buffer{}

	if (info.Mode() & os.ModeCharDevice) == 0 {
		if _, err := io.Copy(&buf, os.Stdin); err != nil {
			return nil, errors.Wrap(err, "failed to read from stdin")
		}
	} else {
		cmd := exec.Command("vipe")
		cmd.Stdin = bytes.NewReader(old)
		cmd.Stdout = &buf
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, errors.Wrap(err, "failed to run editor")
		}
	}

	return buf.Bytes(), nil
}
