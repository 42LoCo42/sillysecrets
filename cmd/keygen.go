package cmd

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:     "keygen <files...>",
	Aliases: []string{"k"},

	Short: "Generate some secret key files",
	Long: `Generate some secret key files.
You can also use basically any file as a key
(e.g. your SSH keys, which will be loaded by default).
The actual secret key is derived from the file data using Argon2id.`,

	Args: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			file, err := os.Create(path)
			if err != nil {
				return errors.Wrap(err, "failed to create key file")
			}
			defer file.Close()

			// 32 bytes = same as internal secret key length
			// (it will be derived from this anyway)
			if _, err := io.CopyN(file, rand.Reader, 32); err != nil {
				return errors.Wrap(err, "failed to generate random data for key file")
			}
		}

		for _, path := range args {
			keyPair := pkg.KeyPair{}
			if err := keyPair.Load(path); err != nil {
				return errors.Wrap(err, "failed to load key file")
			}

			if len(args) > 1 {
				fmt.Printf("%v %v\n", pkg.Encode(keyPair.Public), path)
			} else {
				fmt.Println(pkg.Encode(keyPair.Public))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}
