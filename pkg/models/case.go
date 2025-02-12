package models

import (
	"context"
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/sgrumley/gotex/pkg/runner"
)

var _ Node = (*Case)(nil)

type Case struct {
	Name     string
	Parent   *Function
	Location *ast.KeyValueExpr // the key value of the test case name
}

func (c *Case) GetName() string {
	return c.Name
}

func (c *Case) GetPath() string {
	projectPath := c.Parent.Parent.Parent.Parent.RootDir
	casePath := c.Parent.GetPath() + "/" + c.Name

	return strings.TrimPrefix(casePath, projectPath)
}

func (c *Case) GetChildren() []Node {
	return nil
}

func (c *Case) HasChildren() bool {
	return false
}

func (c *Case) RunTest(ctx context.Context) (*runner.Response, error) {
	function := c.Parent
	file := function.Parent
	pkg := file.Parent
	project := pkg.Parent

	caseName := fmt.Sprintf("%s/%s", function.Name, c.Name)
	path := filepath.Dir(file.Path)
	return runner.RunTest(ctx, runner.TestTypeCase, caseName, path, project.Config)
}
