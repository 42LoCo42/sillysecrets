package cmd

import (
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "sesi",

		Short: "The silliest secret manager! :3",

		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},

		DisableAutoGenTag: true,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
		},
	}

	debug bool

	keyPaths []string
	jsonPath string
	yamlPath string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home := os.Getenv("HOME")

	defaultKeyPaths := []string{
		path.Join(home, ".local/share/sillysecrets"),
	}

	ssh := path.Join(home, ".ssh")
	if entries, err := os.ReadDir(ssh); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "id_") {
				defaultKeyPaths = append(defaultKeyPaths, path.Join(ssh, entry.Name()))
			}
		}
	}

	rootCmd.PersistentFlags().BoolVarP(&debug,
		"debug", "d",
		false,
		"Enable debug logging")

	rootCmd.PersistentFlags().StringArrayVarP(&keyPaths,
		"key", "k",
		defaultKeyPaths,
		"Path to key files")

	rootCmd.PersistentFlags().StringVarP(&jsonPath,
		"json", "j",
		"sesi.json",
		"Path to the JSON storage file")

	rootCmd.PersistentFlags().StringVarP(&yamlPath,
		"yaml", "y",
		"sesi.yaml",
		"Path to the YAML structure file")
}
