package models

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/runner"
)

type Node interface {
	GetName() string
	GetChildren() []Node
	HasChildren() bool
	RunTest() (*runner.Response, error)
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

func (p *Project) RunTest() (*runner.Response, error) {
	path := filepath.Dir(p.RootDir)
	return runner.RunTest(runner.TestTypeProject, "", path, p.Config)
}

// TODO: update to recursive with nodes interface funcs
func (p *Project) FlattenAllNodes() *FlatProject {
	nodes := make(map[string]Node)
	names := make([]string, 0)

	// TODO: update to append names e.g. pkg/file/func/case
	for _, pkg := range p.Packages {
		nodes[pkg.GetName()] = pkg
		names = append(names, pkg.GetName())
		for _, file := range pkg.Files {
			nodes[file.GetName()] = file
			names = append(names, file.GetName())
			for _, function := range file.Functions {
				nodes[function.GetName()] = function
				names = append(names, function.GetName())
				for _, c := range function.Cases {
					nodes[c.GetName()] = c
					names = append(names, c.GetName())
				}
			}
		}
	}
	return &FlatProject{
		Names:   names,
		NodeMap: nodes,
	}
}
