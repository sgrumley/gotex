package finder

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sgrumley/gotex/pkg/config"
	"sgrumley/gotex/pkg/runner"
	"strings"
)

type Node interface {
	GetName() string
	GetChildren() []Node
	HasChildren() bool
	RunTest() (string, error)
}

var _ Node = (*Project)(nil)

type Project struct {
	Config  config.Config
	RootDir string
	Files   []*File
	FileMap map[string]*File
	log     *slog.Logger
}

func (p *Project) GetName() string {
	paths := strings.Split(p.RootDir, "/")

	return paths[len(paths)-1]
}

func (p *Project) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range p.Files {
		children = append(children, child)
	}
	return children
}

func (p *Project) HasChildren() bool {
	if len(p.Files) > 0 {
		return true
	}

	return false
}

func (p *Project) RunTest() (string, error) {
	path := filepath.Dir(p.RootDir)
	return runner.RunTest(runner.TEST_TYPE_PROJECT, "", path, p.Config)
}

func InitProject(log *slog.Logger) (*Project, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load a config: %w", err)
	}

	p := &Project{
		Config: cfg,
		log:    log,
	}
	projectRoot, err := FindGoProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %s\n", err)
	}

	files, err := ListTestFilesWithCWD()
	if err != nil {
		return nil, fmt.Errorf("failed finding any test files: %s\n", err)
	}

	// PERF: this could be concurrent
	for _, file := range files {
		log.Info("searching file: ",
			slog.String("file", file.Path),
		)
		file.Functions = make([]*Function, 0)
		file.FunctionMap = make(map[string]*Function)
		file.Parent = p

		err := SearchFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed finding tests: %s\n", err)
		}
	}

	fileMap := make(map[string]*File)
	fileList := make([]*File, 0)
	for _, file := range files {
		fileMap[file.Name] = file
		fileList = append(fileList, file)
	}

	p.RootDir = projectRoot
	p.Files = fileList
	p.FileMap = fileMap

	log.Info("project starting data",
		slog.String("root dir", p.RootDir),
	)

	return p, nil
}

// func (p *Project) TestNameOut() ([]string, map[string]string) {
// 	tests := make([]string, 0)
// 		// PERF: concurrent here
// 		for _, fn := range f.Functions {
// 			for _, c := range fn.Cases {
// 				tc := fmt.Sprintf("%s/%s", fn.Name, c.Name)
// 				tests = append(tests, tc)
// 				testLocation[tc] = filepath.Dir(f.Path)
// 			}
// 		}
// 	}
// 	return tests, testLocation
// }
