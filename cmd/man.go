package cmd

import (
	"os"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manCmd = &cobra.Command{
	Use:    "man",
	Hidden: true,

	Short: "Generate the manpages",

	Args: cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		path := "man"

		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.Wrap(err, "failed to create manpage directory")
		}

		hdr := &doc.GenManHeader{
			Manual: "sillysecrets",
			Source: "https://github.com/42LoCo42/sillysecrets",
			Title:  "TITLE",
		}

		if err := doc.GenManTree(rootCmd, hdr, path); err != nil {
			return errors.Wrap(err, "failed to generate manpages")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(manCmd)
}
