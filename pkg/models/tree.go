package models

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
)

type DirectoryContent struct {
	name string
	path string
}

func (d DirectoryContent) GetName() string {
	return d.name
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

func (p *Project) PopulateFromPackages(pkgs []*Package) error {
	tree := NewTree(p)
	dirNodes := make(map[string]*NodeTree)
	dirNodes[p.GetName()] = tree.RootNode

	for _, pkg := range pkgs {
		if err := processPackage(p, pkg, tree, dirNodes); err != nil {
			return fmt.Errorf("error processing package %s: %w", pkg.Name, err)
		}
	}

	p.Tree = tree
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

// processPackage handles the creation of package nodes and their files
func processPackage(p *Project, pkg *Package, tree *Tree, dirNodes map[string]*NodeTree) error {
	components, err := getPathComponents(p.RootDir, pkg.Path)
	if err != nil {
		return err
	}

	parentNode, level := createDirectoryStructure(p.RootDir, components, tree.RootNode, dirNodes)
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

	if err := processPackageFiles(pkg, pkgNode, p.log); err != nil {
		return err
	}

	return nil
}

// getPathComponents returns the path components relative to the root directory
func getPathComponents(rootDir, pkgPath string) ([]string, error) {
	relPath, err := filepath.Rel(rootDir, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %w", err)
	}
	return strings.Split(filepath.Clean(relPath), string(filepath.Separator)), nil
}

// createDirectoryStructure creates intermediate directory nodes and returns the last parent node
func createDirectoryStructure(rootPath string, components []string, rootNode *NodeTree, dirNodes map[string]*NodeTree) (*NodeTree, int) {
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
			Data: DirectoryContent{
				name: component,
				path: currentPath,
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

func processPackageFiles(pkg *Package, pkgNode *NodeTree, logger *slog.Logger) error {
	for _, file := range pkg.Files {
		if err := processFile(file, pkg, pkgNode, logger); err != nil {
			return fmt.Errorf("error processing file %s: %w", file.Path, err)
		}
	}
	return nil
}

// processFile handles individual file processing and node creation
func processFile(file *File, pkg *Package, pkgNode *NodeTree, logger *slog.Logger) error {
	logger.Info("searching file:", slog.String("file", file.Path))

	file.Functions = make([]*Function, 0)
	file.FunctionMap = make(map[string]*Function)
	file.Parent = pkg

	fileNode := &NodeTree{
		Level:  pkgNode.Level + 1,
		Data:   file,
		Type:   NODE_TYPE_FILE,
		Parent: pkgNode,
	}
	pkgNode.Children = append(pkgNode.Children, fileNode)

	if err := SearchFile(file, logger, fileNode); err != nil {
		return fmt.Errorf("failed finding tests: %w", err)
	}

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
