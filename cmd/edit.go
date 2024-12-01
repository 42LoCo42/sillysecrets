package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the secret specified by <group>.<secret>",

	ValidArgsFunction: validSecretArgsFunction,

	Run: func(_ *cobra.Command, args []string) {
		s := loadSecret(args)

		cmd := exec.Command("vipe")
		new := bytes.Buffer{}
		cmd.Stdin = bytes.NewReader(s.Value)
		cmd.Stdout = &new
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(errors.Wrap(err, "vipe failed"))
		}

		enc, err := sillysecrets.Encrypt(new.Bytes(), s.GroupName, groups())
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not encrypt data"))
		}

		s.Group.Secrets[s.SecretName] = enc
		_groups[s.GroupName] = s.Group

		if err := sillysecrets.Save(file, groups()); err != nil {
			log.Fatal(errors.Wrap(err, "could not save groups"))
		}
	}}

func init() {
	rootCmd.AddCommand(editCmd)
}
