package models

import (
	"strings"

	"github.com/sgrumley/gotex/pkgv2/runner"
)

var _ Node = (*Package)(nil)

type Package struct {
	Name   string
	Path   string
	Parent *Project
	Files  []*File
	// maybe file map?
}

func (p *Package) GetName() string {
	namePath := strings.Split(p.Name, "/")
	return namePath[len(namePath)-1]
}

func (p *Package) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range p.Files {
		children = append(children, child)
	}
	return children
}

func (p *Package) HasChildren() bool {
	if len(p.Files) > 0 {
		return true
	}
	return false
}

func (p *Package) RunTest() (*runner.Response, error) {
	project := p.Parent

	return runner.RunTest(runner.TestTypePackage, p.Name, p.Path, project.Config)
}
