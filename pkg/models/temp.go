package models

import (
	"go/ast"
	"log/slog"
	"strings"
)

// processTestFunction handles the processing of a potential test function call
func processTestFunction(file *File, rootNode *ast.File, callExpr *ast.CallExpr, log *slog.Logger, fileNode *NodeTree) {
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
func processFunctionAndCases(file *File, fn *ast.FuncDecl, subtestName string, log *slog.Logger, fileNode *NodeTree) {
	log.Debug("test case found",
		slog.String("case name", subtestName),
		slog.String("function name", fn.Name.Name),
		slog.String("file name", file.Name),
	)

	caseName := extractCaseName(subtestName, log)
	function := &Function{
		Name:    fn.Name.Name,
		decl:    fn,
		VarName: subtestName,
		Parent:  file,
	}
	fnNode := &NodeTree{
		Level:  fileNode.Level + 1,
		Data:   function,
		Type:   NODE_TYPE_FUNCTION,
		Parent: fileNode,
	}
	fileNode.Children = append(fileNode.Children, fnNode)

	cases := findValuesOfIndexedField(fn, caseName)
	populateFunctionCases(function, cases)
	file.FunctionMap[function.Name] = function
	file.Functions = append(file.Functions, function)

	for _, tc := range cases {
		caseNode := &NodeTree{
			Level:  fileNode.Level + 2,
			Data:   tc,
			Type:   NODE_TYPE_CASE,
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
func populateFunctionCases(function *Function, cases []*Case) {
	caseMap := make(map[string]*Case)
	for i := range cases {
		cases[i].Parent = function
		caseMap[cases[i].Name] = cases[i]
	}

	function.Cases = cases
	function.CaseMap = caseMap
}
