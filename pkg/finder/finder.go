package finder

import (
	"fmt"
	"go/ast"
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
	return "", nil
}

// type Packages struct {
// TODO: this needs to be linked to the files and be included in project
// }

var _ Node = (*File)(nil)

type File struct {
	Name        string
	Path        string
	Functions   []*Function
	FunctionMap map[string]*Function
	Parent      *Project
}

func (f *File) GetName() string {
	path := fmt.Sprintf("%s/%s", f.Path, f.Name)
	paths := strings.Split(path, "/")
	nodeName := fmt.Sprintf("%s/%s", paths[len(paths)-1], paths[len(paths)-2])

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
	return "", nil
}

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
	return "", nil
}

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

// TODO: cfg should be set within project init
func (c *Case) RunTest() (string, error) {
	function := c.Parent
	file := function.Parent

	caseName := fmt.Sprintf("%s/%s", function.Name, c.Name)
	path := filepath.Dir(file.Path)
	cfg := config.Config{}
	return runner.RunTest(caseName, path, cfg)
}

func InitProject() *Project {
	p := &Project{}
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

		err := ListAll(file)
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

func (p *Project) TestNameOut() ([]string, map[string]string) {
	tests := make([]string, 0)
	testLocation := make(map[string]string)
	for _, f := range p.Files {
		// PERF: concurrent here
		for _, fn := range f.Functions {
			for _, c := range fn.Cases {
				tc := fmt.Sprintf("%s/%s", fn.Name, c.Name)
				tests = append(tests, tc)
				testLocation[tc] = filepath.Dir(f.Path)
			}
		}
	}
	return tests, testLocation
}
