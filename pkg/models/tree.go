package models

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

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

// type Content interface {
// 	GetName() string
// 	// GetData() will return the needed from within the concrete type
// }

type Tree struct {
	RootNode    *NodeTree
	TotalLevels int
}

type NodeTree struct {
	Level    int
	Data     Node
	Type     NodeType
	Children []*NodeTree
	Parent   *NodeTree
}

func (t *Tree) Traverse(f func(*NodeTree) error) error {
	if t == nil {
		return fmt.Errorf("nil node when traversing")
	}
	if t.RootNode == nil {
		return nil
	}
	return traverse(t.RootNode, f)
}

func traverse(node *NodeTree, f func(*NodeTree) error) error {
	if node == nil {
		return fmt.Errorf("nil node when traversing")
	}
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

func GenerateTree(p *Project) error {
	tree := NewTree(p)
	dirNodes := make(map[string]*NodeTree)
	dirNodes[p.GetName()] = tree.RootNode

	for _, pkg := range p.Packages {
		if err := processPackage(p, pkg, tree, dirNodes); err != nil {
			return fmt.Errorf("failed generating tree: %w", err)
		}
	}
	p.Tree = tree

	return nil
}

func processPackage(p *Project, pkg *Package, tree *Tree, dirNodes map[string]*NodeTree) error {
	relPath, err := filepath.Rel(p.RootDir, pkg.Path)
	if err != nil {
		return fmt.Errorf("error getting relative path: %w", err)
	}
	components := strings.Split(filepath.Clean(relPath), string(filepath.Separator))

	parentNode, level := createDirectoryNodes(p.RootDir, components, tree.RootNode, dirNodes)
	if level+1 > tree.TotalLevels {
		tree.TotalLevels = level + 1
	}

	pkgNode := &NodeTree{
		Level:  level + 1,
		Data:   pkg,
		Type:   NODE_TYPE_PACKAGE,
		Parent: parentNode,
	}
	parentNode.Children = append(parentNode.Children, pkgNode)

	for _, file := range pkg.Files {
		if err := processFile(file, pkgNode); err != nil {
			return fmt.Errorf("error processing file %s: %w", file.Path, err)
		}
	}

	return nil
}

func createDirectoryNodes(rootPath string, components []string, rootNode *NodeTree, dirNodes map[string]*NodeTree) (*NodeTree, int) {
	currentPath := rootPath
	parentNode := rootNode
	maxLevel := 0

	// Process all components except the last one (which will be the package)
	for i, component := range components[:len(components)-1] {
		currentPath = filepath.Join(currentPath, component)
		level := i + 1
		maxLevel = level

		if node, exists := dirNodes[currentPath]; exists {
			parentNode = node
			continue
		}

		newDirNode := &NodeTree{
			Level: level,
			Data: &DirectoryContent{
				Name: component,
				Path: currentPath,
			},
			Type:   NODE_TYPE_DIRECTORY,
			Parent: parentNode,
		}
		parentNode.Children = append(parentNode.Children, newDirNode)

		dirNodes[currentPath] = newDirNode
		parentNode = newDirNode
	}

	return parentNode, maxLevel
}

func processFile(file *File, pkgNode *NodeTree) error {
	fileNode := &NodeTree{
		Level:  pkgNode.Level + 1,
		Data:   file,
		Type:   NODE_TYPE_FILE,
		Parent: pkgNode,
	}
	pkgNode.Children = append(pkgNode.Children, fileNode)

	for _, fn := range file.Functions {
		if err := processFunction(fn, fileNode); err != nil {
			return fmt.Errorf("failed processing function")
		}
	}

	return nil
}

func processFunction(fn *Function, fileNode *NodeTree) error {
	fnNode := &NodeTree{
		Level:  fileNode.Level + 1,
		Data:   fn,
		Type:   NODE_TYPE_FUNCTION,
		Parent: fileNode,
	}

	fileNode.Children = append(fileNode.Children, fnNode)
	for _, tc := range fn.Cases {
		if err := processCase(tc, fnNode); err != nil {
			return fmt.Errorf("failed processing function")
		}
	}
	return nil
}

func processCase(tc *Case, fnNode *NodeTree) error {
	caseNode := &NodeTree{
		Level:  fnNode.Level + 1,
		Data:   tc,
		Type:   NODE_TYPE_CASE,
		Parent: fnNode,
	}

	fnNode.Children = append(fnNode.Children, caseNode)
	return nil
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
