package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the secret specified by <group>.<secret>",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveError
		}

		return secrets().ToSlice(), cobra.ShellCompDirectiveNoFileComp
	}, Run: func(_ *cobra.Command, args []string) {
		name := args[0]
		parts := strings.Split(name, ".")
		if len(parts) != 2 {
			log.Fatalf("invalid secret %v: must be in <group>.<secret> format", name)
		}

		groupName := strings.TrimSpace(parts[0])
		secretName := strings.TrimSpace(parts[1])
		if groupName == "" || secretName == "" {
			log.Fatalf("invalid secret %v: must be in <group>.<secret> format", name)
		}

		group, ok := groups()[groupName]
		if !ok {
			log.Fatalf("group %v not found", groupName)
		}

		enc, ok := group.Secrets[secretName]
		old := []byte{}
		if ok {
			var err error
			old, err = sillysecrets.Decrypt(enc, ids())
			if err != nil {
				log.Fatal(errors.Wrap(err, "could not decrypt data"))
			}
		}

		cmd := exec.Command("vipe")
		new := bytes.Buffer{}
		cmd.Stdin = bytes.NewReader(old)
		cmd.Stdout = &new
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(errors.Wrap(err, "vipe failed"))
		}

		var err error
		enc, err = sillysecrets.Encrypt(new.Bytes(), groupName, groups())
		if err != nil {
			log.Fatal(errors.Wrap(err, "could not encrypt data"))
		}

		if group.Secrets == nil {
			group.Secrets = map[string]string{}
		}

		group.Secrets[secretName] = enc
		_groups[groupName] = group

		if err := sillysecrets.Save(file, groups()); err != nil {
			log.Fatal(errors.Wrap(err, "could not save groups"))
		}
	}}

func init() {
	rootCmd.AddCommand(editCmd)
}
