package finder

import (
	"fmt"
	"log"
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

func InitProject() *Project {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to load a config", err)
	}

	p := &Project{
		Config: cfg,
	}
	projectRoot, err := FindGoProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %s\n", err)
	}

	files, err := ListTestFilesWithCWD()
	if err != nil {
		log.Fatalf("failed finding any test files: %s\n", err)
	}

	// PERF: this could be concurrent
	for _, file := range files {
		fmt.Printf("searching file: %s\n", file.Path)
		file.Functions = make([]*Function, 0)
		file.FunctionMap = make(map[string]*Function)
		file.Parent = p

		err := SearchFile(file)
		if err != nil {
			log.Fatalf("failed finding tests: %s\n", err)
		}
	}

	// map files to map[name]file
	fileMap := make(map[string]*File)
	fileList := make([]*File, 0)
	for _, file := range files {
		fileMap[file.Name] = file
		fileList = append(fileList, file)
	}

	p.RootDir = projectRoot
	p.Files = fileList
	p.FileMap = fileMap

	return p
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
