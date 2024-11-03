package finder

import (
	"go/ast"
	"path/filepath"
	"sgrumley/gotex/pkg/runner"
)

var _ Node = (*Function)(nil)

type Function struct {
	Name    string
	Cases   []*Case
	CaseMap map[string]*Case
	Parent  *File
	// meta data that may be helpful
	VarName string
	decl    *ast.FuncDecl
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

func (f *Function) RunTest() (string, error) {
	file := f.Parent
	project := file.Parent

	path := filepath.Dir(file.Path)
	return runner.RunTest(runner.TEST_TYPE_FUNCTION, f.Name, path, project.Config)
}
