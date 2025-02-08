package models

import (
	"fmt"
	"go/ast"
	"path/filepath"

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

func (c *Case) GetChildren() []Node {
	return nil
}

func (c *Case) HasChildren() bool {
	return false
}

func (c *Case) RunTest() (*runner.Response, error) {
	function := c.Parent
	file := function.Parent
	pkg := file.Parent
	project := pkg.Parent

	caseName := fmt.Sprintf("%s/%s", function.Name, c.Name)
	path := filepath.Dir(file.Path)
	return runner.RunTest(runner.TestTypeCase, caseName, path, project.Config)
}
