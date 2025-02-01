package models

import (
	"fmt"
	"log"
	"strings"
)

type DirectoryContent struct {
	Name string
	Path string
}

func (d DirectoryContent) GetName() string {
	return d.Name
}

// NOTE: This tree should hold the structure of the data extracted from the project with the levels preserved as the file tree paths
// This should probably replace the project struct
// while implementing this, adding concurrency for branches would be nice
// func (t *Tree) Find(string) {}

type NodeType string

const (
	NODE_TYPE_PROJECT   NodeType = "PROJECT"
	NODE_TYPE_DIRECTORY NodeType = "DIRECTORY"
	NODE_TYPE_PACKAGE   NodeType = "PACKAGE"
	NODE_TYPE_FILE      NodeType = "FILE"
	NODE_TYPE_FUNCTION  NodeType = "FUNCTION"
	NODE_TYPE_CASE      NodeType = "CASE"
)

type Content interface {
	GetName() string
	// GetData() will return the needed from within the concrete type
}

type Tree struct {
	RootNode    *NodeTree
	TotalLevels int
}

type NodeTree struct {
	Level    int
	Data     Content
	Type     NodeType
	Children []*NodeTree
	Parent   *NodeTree
}

func (t *Tree) Traverse(f func(*NodeTree) error) error {
	if t.RootNode == nil {
		return nil
	}
	return traverse(t.RootNode, f)
}

func traverse(node *NodeTree, f func(*NodeTree) error) error {
	if err := f(node); err != nil {
		return err
	}

	for _, child := range node.Children {
		if err := traverse(child, f); err != nil {
			return err
		}
	}

	return nil
}

// NewTree creates and initializes a new tree with the project as root
func NewTree(p *Project) *Tree {
	return &Tree{
		RootNode: &NodeTree{
			Level: 0,
			Data:  p,
			Type:  NODE_TYPE_PROJECT,
		},
	}
}

func (t *Tree) Print() {
	err := t.Traverse(func(node *NodeTree) error {
		indent := strings.Repeat("  ", node.Level)
		prefix := "├── "
		if node.Parent != nil && node == node.Parent.Children[len(node.Parent.Children)-1] {
			prefix = "└── "
		}

		if node.Level == 0 {
			fmt.Printf("%s/\n", node.Data.GetName())
		} else if node.Type == NODE_TYPE_PACKAGE {
			fmt.Printf("%s%s%s (%s)\n", indent, prefix, node.Data.GetName(), node.Data.GetName())
		} else {
			fmt.Printf("%s%s%s/\n", indent, prefix, node.Data.GetName())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
