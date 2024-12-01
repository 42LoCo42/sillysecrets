package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var file string
var idPaths []string

var rootCmd = &cobra.Command{
	Use: "sesi",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&file, "file", "f",
		"sesi.yaml",
		"Path to the sesi storage file")

	rootCmd.PersistentFlags().StringArrayVarP(
		&idPaths, "identity", "i",
		[]string{path.Join(os.Getenv("HOME"), ".ssh")},
		"Use the specified identity(s) for decryption")
}
