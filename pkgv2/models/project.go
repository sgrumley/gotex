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
	if len(p.Packages) > 0 {
		return true
	}

	return false
}

func (p *Project) RunTest() (*runner.Response, error) {
	path := filepath.Dir(p.RootDir)
	return runner.RunTest(runner.TestTypeProject, "", path, p.Config)
}
