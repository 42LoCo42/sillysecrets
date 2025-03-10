package pkg

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/deckarep/golang-set/v2"
	"github.com/go-faster/errors"
)

// An item of the storage tree
type Node struct {
	Name string `yaml:"-"` // Local name of this node. Do not generally use except when validating it.
	Path string `yaml:"-"` // Full path to this node. Use this when referring to it.

	Key  string `yaml:":key,omitempty"`  // A key with which to encrypt secrets of this node
	Keys *Set   `yaml:":keys,omitempty"` // The full set of keys with which to encrypt secrets of this node

	Import *Set `yaml:":import,omitempty"` // List of node names to import secrets from
	Export *Set `yaml:":export,omitempty"` // List of node names to export our own secrets to

	Children Tree  `yaml:",inline"`
	Parent   *Node `yaml:"-"`
}

// cache for completed import/export operations
var done = map[string]struct{}{}

type NodeTraverseF func(this *Node) error

// Call the traversal function for this node and all of its parents
func (node *Node) Traverse(f NodeTraverseF) error {
	if err := f(node); err != nil {
		return err
	}

	if node.Parent != nil {
		return node.Parent.Traverse(f)
	}

	return nil
}

// Add a key to this node
func (node *Node) AddKey(key string) {
	if node.Keys == nil {
		node.Keys = &Set{mapset.NewSet[string]()}
	}

	node.Keys.Add(key)
}

// Recursively enumerates all keys with which to encrypt secrets of this node
func (node *Node) AllKeys() (keys Set) {
	keys = Set{mapset.NewSet[string]()}

	node.Traverse(func(this *Node) error {
		if this.Keys != nil {
			keys = Set{keys.Union(this.Keys.Set)}
		}
		return nil
	})

	return keys
}

// Does this node have an export list defined?
func (node *Node) IsExport() bool {
	return node.Export != nil && node.Export.Cardinality() > 0
}

// Export secrets into the given node
func (this *Node) DoExport(root *Tree) error {
	for _, thatN := range this.Export.ToSlice() {
		that, err := root.Get(thatN)
		if err != nil {
			return errors.Wrapf(err, "`%v`: export", this.Path)
		}

		slog.Debug("processing export",
			slog.String("this", this.Path),
			slog.String("that", that.Path))

		if that.Import == nil {
			that.Import = &Set{mapset.NewSet[string]()}
		}

		that.Import.Add(this.Path)
	}

	return nil
}

// Does this node have an import list defined?
func (node *Node) IsImport() bool {
	return node.Import != nil && node.Import.Cardinality() > 0
}

// Import secrets using the given keys
func (this *Node) DoImport(level int, keys *Set, root *Tree) error {
	if keys == nil || keys.Cardinality() <= 0 {
		return errors.Errorf("`%v` imports secrets, but has no keys!", this.Path)
	}

	for _, thatN := range this.Import.ToSlice() {
		that, err := root.Get(thatN)
		if err != nil {
			return errors.Wrapf(err, "`%v`: import", this.Path)
		}

		// process this import
		for _, key := range keys.ToSlice() {
			id := fmt.Sprintf("%v %v %v", this.Path, that.Path, key)
			if _, ok := done[id]; ok {
				continue
			}
			done[id] = struct{}{}

			slog.Debug(strings.Repeat("  ", level)+"processing import",
				slog.String("this", this.Path),
				slog.String("that", that.Path),
				slog.String("key", key))

			that.AddKey(key)
		}

		// process import at target
		if that.IsImport() {
			if err := that.DoImport(level+1, keys, root); err != nil {
				return err
			}
		}

		// scan target downwards for imports
		if err := that.Children.traverse(that, root,
			func(this *Node, root *Tree) error {
				if this.IsImport() {
					if err := this.DoImport(level+1, keys, root); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
			return err
		}
	}

	return nil
}

// A node is a valid secret if it is a leaf of the tree
func (node *Node) IsSecret() bool {
	return len(node.Children) == 0
}
