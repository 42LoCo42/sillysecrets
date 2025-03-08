package cmd

import (
	"os"
	"strings"

	"github.com/42LoCo42/sillysecrets/pkg"
	"github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var nonames bool
var noexpand bool

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:     "tree",
	Aliases: []string{"t"},

	Short: "View the entire tree after validation",
	Long: `View the entire tree after validation.
This is mostly useful for debugging your import and export rules.
All key sets will be expanded to include their parent's keys, unless the flag -x is given.
Then, all keys that occur as a :key property in a node will be replaced
by that node's uppercased name, unless the flag -n is given.`,

	Args: cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		tree := pkg.Tree{}
		if err := tree.Load(yamlPath); err != nil {
			return errors.Wrap(err, "failed to load tree")
		}

		if !noexpand {
			tree.Traverse(func(this *pkg.Node, root *pkg.Tree) error {
				if this.Keys != nil && this.Keys.Cardinality() > 0 {
					keys := this.AllKeys()
					this.Keys = &keys
				}
				return nil
			})
		}

		if !nonames {
			keyNames := map[string]string{}

			// collect names
			tree.Traverse(func(this *pkg.Node, root *pkg.Tree) error {
				if this.Key != "" {
					keyNames[this.Key] = strings.ToUpper(this.Path)
				}
				return nil
			})

			// replace keys with names
			tree.Traverse(func(this *pkg.Node, root *pkg.Tree) error {
				if this.Keys != nil && this.Keys.Cardinality() > 0 {
					newKeys := &pkg.Set{Set: mapset.NewSet[string]()}

					for _, key := range this.Keys.ToSlice() {
						path, ok := keyNames[key]
						if ok {
							newKeys.Add(path)
						} else {
							newKeys.Add(key)
						}
					}

					this.Keys = newKeys
				}
				return nil
			})
		}

		enc := yaml.NewEncoder(os.Stdout)
		enc.SetIndent(2)
		if err := enc.Encode(tree); err != nil {
			return errors.Wrap(err, "failed to encode tree")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)

	treeCmd.Flags().BoolVarP(&nonames,
		"nonames", "n",
		false,
		"Don't replace known keys with names")

	treeCmd.Flags().BoolVarP(&noexpand,
		"noexpand", "x",
		false,
		"Don't expand key sets")
}
