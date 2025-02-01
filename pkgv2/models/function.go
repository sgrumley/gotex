package models

import (
	"go/ast"
	"path/filepath"

	"github.com/sgrumley/gotex/pkgv2/runner"
)

var _ Node = (*Function)(nil)

type Function struct {
	Name    string
	Cases   []*Case
	CaseMap map[string]*Case
	Parent  *File
	// meta data that may be helpful
	VarName string
	Decl    *ast.FuncDecl
}

func (f *Function) GetName() string {
	return f.Name
}

func (f *Function) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range f.Cases {
		children = append(children, child)
	}
	return children
}

func (f *Function) HasChildren() bool {
	if len(f.Cases) > 0 {
		return true
	}

	return false
}

func (f *Function) RunTest() (*runner.Response, error) {
	file := f.Parent
	pkg := file.Parent
	project := pkg.Parent

	path := filepath.Dir(file.Path)
	return runner.RunTest(runner.TestTypeFunction, f.Name, path, project.Config)
}
