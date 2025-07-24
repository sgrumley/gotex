// Package models provides data structures for representing Go projects,
// packages, files, functions, and test cases in a hierarchical tree structure.
//
// The models package is central to the gotex library, defining the core
// data structures used throughout the test analysis and execution system.
// All model types implement the Node interface, enabling uniform tree
// traversal and test execution operations.
//
// The hierarchy is:
//   Project -> Package -> File -> Function -> Case
//
// Each level supports:
//   - Tree navigation (GetName, GetPath, GetChildren, HasChildren)
//   - Test execution (RunTest)
//   - Structured representation of Go test code
//
// Example usage:
//
//	project := &models.Project{
//		RootDir: "/path/to/project",
//		Config:  cfg,
//	}
//	
//	// Navigate the project structure
//	for _, pkg := range project.Packages {
//		fmt.Printf("Package: %s\n", pkg.GetName())
//	}
package models

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/runner"
)

// Node represents any navigable element in the project tree hierarchy.
// It provides a uniform interface for traversing and executing tests
// across projects, packages, files, functions, and test cases.
type Node interface {
	GetName() string                                    // Returns the display name of the node
	GetPath() string                                    // Returns the hierarchical path from root
	GetChildren() []Node                                // Returns child nodes in the tree
	HasChildren() bool                                  // Indicates if the node has children
	RunTest(ctx context.Context) (*runner.Response, error) // Executes tests for this node
}

var _ Node = (*Project)(nil)

// Project represents the root of a Go project's test structure.
// It contains all packages, maintains configuration, and provides
// the entry point for project-wide operations.
type Project struct {
	Config   config.Config // Test execution configuration
	RootDir  string        // Absolute path to the project root
	Packages []*Package    // All packages in the project
	Tree     *Tree         // Hierarchical tree representation
	Log      *slog.Logger  // Logger instance
}

// FlatProject provides a flattened view of all nodes in the project
// for efficient lookup and navigation operations.
type FlatProject struct {
	NodeMap map[string]Node // Map of path -> Node for O(1) lookup
	Names   []string        // Ordered list of all node paths
}

func (p *Project) GetName() string {
	paths := strings.Split(p.RootDir, "/")

	return paths[len(paths)-1]
}

func (p *Project) GetPath() string {
	return "/"
}

func (p *Project) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range p.Packages {
		children = append(children, child)
	}
	return children
}

func (p *Project) HasChildren() bool {
	return len(p.Packages) > 0
}

func (p *Project) RunTest(ctx context.Context) (*runner.Response, error) {
	path := filepath.Dir(p.RootDir)
	return runner.RunTest(ctx, runner.TestTypeProject, "", path, p.Config)
}

// FlattenAllNodes creates a flattened representation of the entire project tree.
// This enables efficient lookups and provides ordered access to all nodes
// across the project hierarchy. The returned FlatProject contains both
// a map for O(1) lookups and an ordered slice for iteration.
func (p *Project) FlattenAllNodes() *FlatProject {
	nodes := make(map[string]Node)
	names := make([]string, 0)

	for _, pkg := range p.Packages {
		nodes[pkg.GetPath()] = pkg
		names = append(names, pkg.GetPath())
		for _, file := range pkg.Files {
			nodes[file.GetPath()] = file
			names = append(names, file.GetPath())
			for _, function := range file.Functions {
				nodes[function.GetPath()] = function
				names = append(names, function.GetPath())
				for _, c := range function.Cases {
					nodes[c.GetPath()] = c
					names = append(names, c.GetPath())
				}
			}
		}
	}
	return &FlatProject{
		Names:   names,
		NodeMap: nodes,
	}
}

func (p *Project) Print() {
	fmt.Println("\n=== Final Package State ===")
	for i, pkg := range p.Packages {
		fmt.Printf("Package[%d] address: %p\n", i, pkg)
		for j, file := range pkg.Files {
			fmt.Printf("  File[%d] address: %p\n", j, file)
			fmt.Printf("    Functions: %d\n", len(file.Functions))
			for k, fn := range file.Functions {
				fmt.Printf("      Function[%d] address: %p, Cases: %d\n", k, fn, len(fn.Cases))
			}
		}
	}
}
