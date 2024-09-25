package finder

import (
	"fmt"
	"go/ast"
	"log"
)

type Project struct {
	RootDir string
	Files   []*File
}

type File struct {
	Name      string
	Path      string
	Functions []Function
}

type Function struct {
	Name    string
	Cases   []Case
	decl    *ast.FuncDecl
	VarName string
}

type Case struct {
	Name string
}

func InitProject() *Project {
	projectRoot, err := FindGoProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %s", err)
	}

	files, err := ListTestFilesWithCWD()
	if err != nil {
		log.Fatalf("failed finding any test files: %s", err)
	}

	// PERF: this could be concurrent
	for _, file := range files {
		fmt.Printf("searching file: %s\n", file.Path)
		err := ListAll(file)
		if err != nil {
			log.Fatalf("failed finding tests: %s", err)
		}
	}

	return &Project{
		RootDir: projectRoot,
		Files:   files,
	}
}

func (p Project) PrettyPrint() {
	fmt.Printf("\n┌── 📂 %s/\n", p.RootDir)
	for i, file := range p.Files {
		isLastFile := i == len(p.Files)-1
		PrettyPrintFile(file, "", isLastFile)
	}
}

// PrettyPrintFile prints the File struct with tree lines
func PrettyPrintFile(f *File, prefix string, isLast bool) {
	fileBranch := "└──" // Last element
	if !isLast {
		fileBranch = "├──" // Intermediate element
	}
	fmt.Printf("%s%s 📝 %s\n", prefix, fileBranch, f.Name)

	// Update the prefix for nested levels
	newPrefix := prefix
	if isLast {
		newPrefix += "    " // Indent for last element
	} else {
		newPrefix += "│   " // Continue the tree line for intermediate elements
	}

	// Print functions inside the file
	for i, fn := range f.Functions {
		isLastFunc := i == len(f.Functions)-1
		PrettyPrintFunction(fn, newPrefix, isLastFunc)
	}
}

// PrettyPrintFunction prints the Function struct with tree lines
func PrettyPrintFunction(fn Function, prefix string, isLast bool) {
	funcBranch := "└──" // Last element
	if !isLast {
		funcBranch = "├──" // Intermediate element
	}
	fmt.Printf("%s%s 🧪 %s\n", prefix, funcBranch, fn.Name)

	// Update the prefix for cases
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	// Print cases inside the function
	for i, c := range fn.Cases {
		isLastCase := i == len(fn.Cases)-1
		PrettyPrintCase(c, newPrefix, isLastCase)
	}
}

// PrettyPrintCase prints the Case struct with tree lines
func PrettyPrintCase(c Case, prefix string, isLast bool) {
	caseBranch := "└──" // Last element
	if !isLast {
		caseBranch = "├──" // Intermediate element
	}
	fmt.Printf("%s%s 💼 %s\n", prefix, caseBranch, c.Name)
}

// Just tabbed space
// // PrettyPrintProject prints the Project struct in a readable format
// func PrettyPrintProject(p Project) {
// 	fmt.Printf("Project Root: %s\n", p.RootDir)
// 	fmt.Println("Files:")
// 	for _, file := range p.Files {
// 		PrettyPrintFile(file, 1)
// 	}
// }

// // PrettyPrintFile prints the File struct with indentation
// func PrettyPrintFile(f *File, indentLevel int) {
// 	indent := strings.Repeat("  ", indentLevel)
// 	fmt.Printf("%sFile Name: %s\n", indent, f.Name)
// 	fmt.Printf("%sFile Path: %s\n", indent, f.Path)
// 	fmt.Printf("%sFunctions:\n", indent)
// 	for _, fn := range f.Functions {
// 		PrettyPrintFunction(fn, indentLevel+1)
// 	}
// }

// // PrettyPrintFunction prints the Function struct with indentation
// func PrettyPrintFunction(fn Function, indentLevel int) {
// 	indent := strings.Repeat("  ", indentLevel)
// 	fmt.Printf("%sFunction Name: %s\n", indent, fn.Name)
// 	fmt.Printf("%sVariable Name: %s\n", indent, fn.VarName)
// 	fmt.Printf("%sCases:\n", indent)
// 	for _, c := range fn.Cases {
// 		PrettyPrintCase(c, indentLevel+1)
// 	}
// }

// // PrettyPrintCase prints the Case struct with indentation
// func PrettyPrintCase(c Case, indentLevel int) {
// 	indent := strings.Repeat("  ", indentLevel)
// 	fmt.Printf("%sCase Name: %s\n", indent, c.Name)
// }
