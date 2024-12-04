package cmd

import (
	"log"
	"os"
	"path"

	sillysecrets "github.com/42LoCo42/sillysecrets/pkg"
	set "github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
)

var decryptallCmd = &cobra.Command{
	Use:   "decryptall",
	Short: "Decrypt all secrets accessible by a given group into a folder",

	Args: cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return set.NewSetFromMapKeys(groups()).ToSlice(), cobra.ShellCompDirectiveNoFileComp
		case 1:
			return nil, cobra.ShellCompDirectiveFilterDirs
		default:
			return nil, cobra.ShellCompDirectiveError
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]
		dir := args[1]

		secrets, err := sillysecrets.CollectSecrets(group, groups())
		if err != nil {
			log.Fatal(errors.Wrapf(err, "could not collect secrets for group %v", group))
		}

		for _, name := range secrets.ToSlice() {
			secret, err := loadSecret(name)
			if err != nil {
				log.Print("WARN: ", errors.Wrapf(err, "could not load secret %v", name))
				continue
			}

			outpath := path.Join(dir, secret.GroupName, secret.SecretName)
			os.MkdirAll(path.Dir(outpath), 0700)

			if err := os.WriteFile(outpath, secret.Value, 0400); err != nil {
				log.Print("WARN: ", errors.Wrapf(err,
					"could not write secret %v to output file %v",
					secret.RawName, outpath))
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
