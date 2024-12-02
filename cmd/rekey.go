package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rekeyCmd = &cobra.Command{
	Use:   "rekey",
	Short: "Try to rekey all secrets",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveError
	},

	Run: func(cmd *cobra.Command, args []string) {
		for _, name := range secrets().ToSlice() {
			secret, err := loadSecret(name)
			if err != nil {
				log.Printf("WARN: could not rekey %v: %v", name, err)
				continue
			}

			if err := saveSecret(secret); err != nil {
				log.Printf("WARN: could not save rekeyed secret %v: %v", name, err)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rekeyCmd)
}
