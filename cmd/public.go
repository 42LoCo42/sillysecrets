package cmd

import (
	"fmt"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var publicCmd = &cobra.Command{
	Use:     "public <files...>",
	Aliases: []string{"p"},

	Short: "Print public keys corresponding to some secret key files",
	Long: `Print public keys corresponding to some secret key files.
If only a single file is given, print only the public key.
Otherwise, print the public keys together with their respective path,
separated by a space.`,

	Args: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
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
	rootCmd.AddCommand(publicCmd)
}
