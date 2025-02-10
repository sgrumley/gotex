package models

import (
	"context"
	"go/ast"
	"path/filepath"

	"github.com/sgrumley/gotex/pkg/runner"
)

var _ Node = (*Function)(nil)

type Function struct {
	Name    string
	Cases   []*Case
	CaseMap map[string]*Case
	Parent  *File
	// meta data that may be helpful
	VarName          string
	TestFunctionNode *ast.FuncDecl // reference to the function AST node
	RunCallNode      *ast.CallExpr // reference to the function call `t.Run()`
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
	return len(f.Cases) > 0
}

func (f *Function) RunTest(ctx context.Context) (*runner.Response, error) {
	file := f.Parent
	pkg := file.Parent
	project := pkg.Parent

	path := filepath.Dir(file.Path)
	return runner.RunTest(ctx, runner.TestTypeFunction, f.Name, path, project.Config)
}
