package finder

import (
	"fmt"
	"path/filepath"
	"sgrumley/gotex/pkg/runner"
)

var _ Node = (*Case)(nil)

type Case struct {
	Name   string
	Parent *Function
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

func (c *Case) RunTest() (string, error) {
	function := c.Parent
	file := function.Parent
	pkg := file.Parent
	project := pkg.Parent

	caseName := fmt.Sprintf("%s/%s", function.Name, c.Name)
	path := filepath.Dir(file.Path)
	return runner.RunTest(runner.TEST_TYPE_CASE, caseName, path, project.Config)
}
