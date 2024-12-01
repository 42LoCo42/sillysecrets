package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt the secret specified by <group>.<secret>",

	ValidArgsFunction: validSecretArgsFunction,

	Run: func(_ *cobra.Command, args []string) {
		s := loadSecret(args)

		if len(s.Value) == 0 {
			log.Fatalf("secret %v is not defined", s.RawName)
		}

		fmt.Print(string(s.Value))
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
