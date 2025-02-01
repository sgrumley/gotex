package scanner

import (
	"fmt"
	"go/ast"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/sgrumley/gotex/pkgv2/models"
)

// processTestFunction handles the processing of a potential test function call
func processTestFunction(file *models.File, rootNode *ast.File, callExpr *ast.CallExpr, log *slog.Logger, fileNode *models.NodeTree) {
	_, ok := isTestRunCall(callExpr)
	if !ok {
		return
	}

	fn := findEnclosingFunction(rootNode, callExpr)
	if fn == nil {
		return
	}

	// TODO: move out of the AST.walk
	// subtestName is the name of the struct element provided to t.Run(tc.name,...)
	subtestName := exprToString(callExpr.Args[0])
	processFunctionAndCases(file, fn, subtestName, log, fileNode)
}

// processFunctionAndCases creates function and test case data structures
func processFunctionAndCases(file *models.File, fn *ast.FuncDecl, subtestName string, log *slog.Logger, fileNode *models.NodeTree) {
	log.Debug("test case found",
		slog.String("case name", subtestName),
		slog.String("function name", fn.Name.Name),
		slog.String("file name", file.Name),
	)

	caseName := extractCaseName(subtestName, log)
	function := &models.Function{
		Name:    fn.Name.Name,
		Decl:    fn,
		VarName: subtestName,
		Parent:  file,
	}
	fnNode := &models.NodeTree{
		Level:  fileNode.Level + 1,
		Data:   function,
		Type:   models.NODE_TYPE_FUNCTION,
		Parent: fileNode,
	}
	fileNode.Children = append(fileNode.Children, fnNode)

	cases := findValuesOfIndexedField(fn, caseName)
	populateFunctionCases(function, cases)
	file.FunctionMap[function.Name] = function
	file.Functions = append(file.Functions, function)

	for _, tc := range cases {
		caseNode := &models.NodeTree{
			Level:  fileNode.Level + 2,
			Data:   tc,
			Type:   models.NODE_TYPE_CASE,
			Parent: fnNode,
		}

		fnNode.Children = append(fnNode.Children, caseNode)
	}
}

// extractCaseName gets the case field name from the subtest name
func extractCaseName(subtestName string, log *slog.Logger) string {
	caseName := "name"
	subtestNameSplit := strings.Split(subtestName, ".")

	if len(subtestNameSplit) == 2 {
		caseName = subtestNameSplit[1]
	} else {
		log.Error("failed identifying struct.name, defaulting to tc.name",
			slog.String("name", subtestName))
	}

	return caseName
}

// populateFunctionCases populates the cases for a function
func populateFunctionCases(function *models.Function, cases []*models.Case) {
	caseMap := make(map[string]*models.Case)
	for i := range cases {
		cases[i].Parent = function
		caseMap[cases[i].Name] = cases[i]
	}

	function.Cases = cases
	function.CaseMap = caseMap
}

// processPackage handles the creation of package nodes and their files
func processPackage(p *models.Project, pkg *models.Package, tree *models.Tree, dirNodes map[string]*models.NodeTree) error {
	components, err := getPathComponents(p.RootDir, pkg.Path)
	if err != nil {
		return err
	}

	parentNode, level := createDirectoryStructure(p.RootDir, components, tree.RootNode, dirNodes)
	if level+1 > tree.TotalLevels {
		tree.TotalLevels = level + 1
	}

	pkgNode := &models.NodeTree{
		Level:  level + 1,
		Data:   pkg,
		Type:   models.NODE_TYPE_PACKAGE,
		Parent: parentNode,
	}
	parentNode.Children = append(parentNode.Children, pkgNode)

	if err := processPackageFiles(pkg, pkgNode, p.Log); err != nil {
		return err
	}

	return nil
}

// getPathComponents returns the path components relative to the root directory
func getPathComponents(rootDir, pkgPath string) ([]string, error) {
	relPath, err := filepath.Rel(rootDir, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path: %w", err)
	}
	return strings.Split(filepath.Clean(relPath), string(filepath.Separator)), nil
}

// createDirectoryStructure creates intermediate directory nodes and returns the last parent node
func createDirectoryStructure(rootPath string, components []string, rootNode *models.NodeTree, dirNodes map[string]*models.NodeTree) (*models.NodeTree, int) {
	currentPath := rootPath
	parentNode := rootNode
	maxLevel := 0

	// Process all components except the last one (which will be the package)
	for i, component := range components[:len(components)-1] {
		currentPath = filepath.Join(currentPath, component)
		level := i + 1
		maxLevel = level

		if node, exists := dirNodes[currentPath]; exists {
			parentNode = node
			continue
		}

		newDirNode := &models.NodeTree{
			Level: level,
			Data: models.DirectoryContent{
				Name: component,
				Path: currentPath,
			},
			Type:   models.NODE_TYPE_DIRECTORY,
			Parent: parentNode,
		}
		parentNode.Children = append(parentNode.Children, newDirNode)

		dirNodes[currentPath] = newDirNode
		parentNode = newDirNode
	}

	return parentNode, maxLevel
}

func processPackageFiles(pkg *models.Package, pkgNode *models.NodeTree, logger *slog.Logger) error {
	for _, file := range pkg.Files {
		if err := processFile(file, pkg, pkgNode, logger); err != nil {
			return fmt.Errorf("error processing file %s: %w", file.Path, err)
		}
	}
	return nil
}

// processFile handles individual file processing and node creation
func processFile(file *models.File, pkg *models.Package, pkgNode *models.NodeTree, logger *slog.Logger) error {
	logger.Info("searching file:", slog.String("file", file.Path))

	file.Functions = make([]*models.Function, 0)
	file.FunctionMap = make(map[string]*models.Function)
	file.Parent = pkg

	fileNode := &models.NodeTree{
		Level:  pkgNode.Level + 1,
		Data:   file,
		Type:   models.NODE_TYPE_FILE,
		Parent: pkgNode,
	}
	pkgNode.Children = append(pkgNode.Children, fileNode)

	if err := SearchFile(file, logger, fileNode); err != nil {
		return fmt.Errorf("failed finding tests: %w", err)
	}

	return nil
}

func PopulateFromPackages(p *models.Project, pkgs []*models.Package) error {
	tree := models.NewTree(p)
	dirNodes := make(map[string]*models.NodeTree)
	dirNodes[p.GetName()] = tree.RootNode

	for _, pkg := range pkgs {
		if err := processPackage(p, pkg, tree, dirNodes); err != nil {
			return fmt.Errorf("error processing package %s: %w", pkg.Name, err)
		}
	}

	p.Tree = tree
	return nil
}
