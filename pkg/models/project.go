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

type Node interface {
	GetName() string
	GetPath() string
	GetChildren() []Node
	HasChildren() bool
	RunTest(ctx context.Context) (*runner.Response, error)
}

var _ Node = (*Project)(nil)

type Project struct {
	Config   config.Config
	RootDir  string
	Packages []*Package
	Tree     *Tree
	Log      *slog.Logger
}

type FlatProject struct {
	NodeMap map[string]Node
	Names   []string
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

// TODO: update to recursive with nodes interface funcs
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
