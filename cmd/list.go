package cmd

import (
	"fmt"
	"log"

	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets accessible by a given group",

	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return set.NewSetFromMapKeys(groups()).ToSlice(), cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveError
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]

		secrets, err := sillysecrets.CollectSecrets(group, groups())
		if err != nil {
			log.Fatal(errors.Wrapf(err, "could not collect secrets for group %v", group))
		}

		for _, name := range secrets.ToSlice() {
			fmt.Println(name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
