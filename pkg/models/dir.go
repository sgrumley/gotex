package models

import (
	"context"

	"github.com/sgrumley/gotex/pkg/runner"
)

var _ Node = (*DirectoryContent)(nil)

type DirectoryContent struct {
	Name string
	Path string
}

func (d *DirectoryContent) GetName() string {
	return d.Name
}

func (d *DirectoryContent) GetChildren() []Node {
	// TODO: what effects does this cause
	return nil
}

func (d *DirectoryContent) HasChildren() bool {
	return true
}

func (d *DirectoryContent) RunTest(ctx context.Context) (*runner.Response, error) {
	return &runner.Response{
		TestType:       runner.TestTypeFile,
		Result:         "Test file not supported",
		Output:         "Test file not supported",
		Error:          "Test file not supported",
		ExternalOutput: "Test file not supported",
		ExternalError:  "Test file not supported",
	}, nil
}
