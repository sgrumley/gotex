package finder

import (
	"fmt"
	"go/ast"
	"log"
	"path/filepath"
)

type Project struct {
	RootDir string
	// Files   []*File
	Files map[string]*File
}

// type Packages struct {
// TODO: this needs to be linked to the files and be included in project
// }

type File struct {
	Name string
	Path string
	// Functions []*Function
	Functions map[string]*Function
}

type Function struct {
	Name string
	// Cases []Case
	Cases map[string]*Case
	// meta data that may be helpful
	VarName string
	decl    *ast.FuncDecl
}

type Case struct {
	Name string
}

func InitProject() *Project {
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
		file.Functions = make(map[string]*Function)
		err := ListAll(file)
		if err != nil {
			log.Fatalf("failed finding tests: %s\n", err)
		}
	}

	// map files to map[name]file
	fileMap := make(map[string]*File)

	for _, file := range files {
		fileMap[file.Name] = file
	}

	return &Project{
		RootDir: projectRoot,
		Files:   fileMap,
	}
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
