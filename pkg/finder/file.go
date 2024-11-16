package finder

import (
	"fmt"
	"path/filepath"
	"sgrumley/gotex/pkg/runner"
	"strings"
)

var _ Node = (*File)(nil)

type File struct {
	Name        string
	Path        string
	Functions   []*Function
	FunctionMap map[string]*Function
	Parent      *Project
}

func (f *File) GetName() string {
	paths := strings.Split(f.Path, "/")
	nodeName := fmt.Sprintf("%s/%s", paths[len(paths)-2], paths[len(paths)-1])
	return nodeName
}

func (f *File) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range f.Functions {
		children = append(children, child)
	}
	return children
}

func (f *File) HasChildren() bool {
	if len(f.Functions) > 0 {
		return true
	}

	return false
}

func (f *File) RunTest() (string, error) {
	project := f.Parent
	path := filepath.Dir(f.Path)

	return runner.RunTest(runner.TEST_TYPE_FILE, f.Name, path, project.Config)
}
