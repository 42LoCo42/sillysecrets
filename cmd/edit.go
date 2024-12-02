package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the secret specified by <group>.<secret>",

	ValidArgsFunction: validSecretArgsFunction,

	Run: func(_ *cobra.Command, args []string) {
		s, err := loadSecret(args[0])
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("vipe")
		new := bytes.Buffer{}
		cmd.Stdin = bytes.NewReader(s.Value)
		cmd.Stdout = &new
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(errors.Wrap(err, "vipe failed"))
		}

		s.Value = new.Bytes()
		if err := saveSecret(s); err != nil {
			log.Fatal(err)
		}
	}}

func init() {
	rootCmd.AddCommand(editCmd)
}
