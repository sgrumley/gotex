package finder

import (
	"fmt"
	"sgrumley/gotex/pkg/runner"
	"strings"
)

var _ Node = (*File)(nil)

type File struct {
	Name        string
	Path        string
	Functions   []*Function
	FunctionMap map[string]*Function
	Parent      *Package
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

func (f *File) RunTest() (*runner.Response, error) {
	// project := f.Parent.Parent
	// path := filepath.Dir(f.Path)
	//
	// return runner.RunTest(runner.TestTypeFile, f.Name, path, project.Config)
	return &runner.Response{
		TestType:       runner.TestTypeFile,
		Result:         "Test file not supported",
		Output:         "Test file not supported",
		Error:          "Test file not supported",
		ExternalOutput: "Test file not supported",
		ExternalError:  "Test file not supported",
	}, nil
}
