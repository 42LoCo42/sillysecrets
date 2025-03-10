package pkg

import (
	"encoding/base64"

	"github.com/spf13/cobra"
)

var format *base64.Encoding = base64.RawURLEncoding

// Encode data with base64 (URL format, no padding)
func Encode(data []byte) string {
	return format.EncodeToString(data)
}

// Decode data from base64 (URL format, no padding)
func Decode(enc string) ([]byte, error) {
	return format.DecodeString(enc)
}

// Given a cobra command with the flag "yaml"
// specifying the path to a structure/tree file,
// act as a shell completion function
// that returns all node/secret names
func CobraValidSecrets(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveError
	}

	yamlPath, err := cmd.Flags().GetString("yaml")
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	tree := Tree{}
	if err := tree.Load(yamlPath); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return tree.Leaves().ToSlice(), cobra.ShellCompDirectiveNoFileComp
}
