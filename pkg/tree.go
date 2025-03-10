package pkg

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
)

// All node names must match this regex
var NameRGX = regexp.MustCompile("^[a-zA-Z0-9._-]+$")

// A map of node names -> nodes
type Tree map[string]*Node

// Load the tree from a YAML file, validating it
func (tree *Tree) Load(path string) error {
	slog.Debug("loading tree from", slog.String("file", path))

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open tree structure file")
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&tree); err != nil {
		return errors.Wrap(err, "failed to save tree")
	}

	if err := tree.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate tree")
	}

	return nil
}

type TreeTraverseF func(this *Node, root *Tree) error

// Call the traversal function recursively for all child nodes in the tree
func (tree *Tree) Traverse(f TreeTraverseF) error {
	return tree.traverse(&Node{}, tree, f)
}

// Internal traverse implementation, taking a parent and root node
// instead of defaulting to an empty root node
func (tree *Tree) traverse(parent *Node, root *Tree, f TreeTraverseF) error {
	for name, this := range *tree {
		if this == nil {
			this = new(Node)

			if parent.Children == nil {
				parent.Children = Tree{}
			}

			parent.Children[name] = this
		}

		if this.Name == "" {
			this.Name = name
		}

		if this.Path == "" {
			path := name
			if parent.Path != "" {
				path = fmt.Sprintf("%v/%v", parent.Path, path)
			}
			this.Path = path
		}

		if this.Parent == nil {
			this.Parent = parent
		}

		if err := f(this, root); err != nil {
			return err
		}

		if err := this.Children.traverse(this, root, f); err != nil {
			return err
		}
	}

	return nil
}

// Get a node from the tree by its path
func (tree *Tree) Get(path string) (node *Node, err error) {
	node = nil

	tree.Traverse(func(this *Node, root *Tree) error {
		if this.Path == path {
			node = this
		}

		return nil
	})

	if node == nil {
		return nil, errors.Errorf("no such node `%v`", path)
	}

	return node, nil
}

// Validate the tree. This is done in three phases:
// 1. validate all node names & include Key in Keys
// 2. perform all exports (i.e. resolve them to imports)
// 3. perform all imports
func (tree *Tree) Validate() error {
	if err := tree.Traverse(func(this *Node, root *Tree) error {
		if !NameRGX.MatchString(this.Name) {
			return errors.Errorf(
				"`%v`: invalid node name `%v`",
				this.Parent.Path, this.Name)
		}

		if this.Key != "" {
			this.AddKey(this.Key)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "tree validation: phase 1 failed")
	}

	if err := tree.Traverse(func(this *Node, root *Tree) error {
		if this.IsExport() {
			if err := this.DoExport(root); err != nil {
				return errors.Wrapf(err, "export failed for `%v`", this.Path)
			}
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "tree validation: phase 2 failed")
	}

	if err := tree.Traverse(func(this *Node, root *Tree) error {
		if this.IsImport() {
			if err := this.DoImport(0, this.Keys, root); err != nil {
				return errors.Wrapf(err, "import failed for `%v`", this.Path)
			}
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "tree validation: phase 3 failed")
	}

	return nil
}

// Return all leaf paths. These elements are valid secrets.
func (tree *Tree) Leaves() Set {
	names := Set{mapset.NewSet[string]()}

	tree.Traverse(func(this *Node, root *Tree) error {
		if this.IsSecret() {
			names.Add(this.Path)
		}
		return nil
	})

	return names
}
