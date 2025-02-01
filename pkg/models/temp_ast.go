package models

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strings"
)

// SearchFile analyzes a Go file for test functions and their test cases
func SearchFile(file *File, log *slog.Logger, fileNode *NodeTree) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file.Path, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", file.Path, err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			processTestFunction(file, node, callExpr, log, fileNode)
		}
		return true
	})

	return nil
}

// isTestRunCall assumes all test functions contain a call to t.Run()
// NOTE: This will break if t.Run() is not used
// TODO: expression shouldn't be hard coded t, it should be determined by the parameters into the function
func isTestRunCall(callExpr *ast.CallExpr) (*ast.SelectorExpr, bool) {
	// a selector expression is the syntax used to access fields or methods of a struct
	// checking that the selector is `Run`, the expression is then checked for `t`
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Run" {
		return nil, false
	}

	if ex := exprToString(selExpr.X); ex != "t" {
		return nil, false
	}

	return selExpr, true
}

// findEnclosingFunction traverses the AST upwards to find the function that encloses the node
// TODO: this is disgusting and needs cleaning
func findEnclosingFunction(node ast.Node, n ast.Node) *ast.FuncDecl {
	// Walk the AST to find the enclosing function
	var fn *ast.FuncDecl
	ast.Inspect(node, func(x ast.Node) bool {
		if f, ok := x.(*ast.FuncDecl); ok {
			if f.Body != nil && f.Body.Pos() <= n.Pos() && f.Body.End() >= n.End() {
				fn = f
				return false // Stop traversing
			}
		}
		return true
	})
	return fn
}

// TODO: this is disgusting and needs cleaning
func findValuesOfIndexedField(fn *ast.FuncDecl, fieldName string) []*Case {
	var cases []*Case

	// findValuesOfIndexedField looks for the value of a field in an array or slice (e.g. tc[i].name)
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// We're looking for composite literals (array/slice initialization) or assignments
		if compLit, ok := n.(*ast.CompositeLit); ok {
			for _, elt := range compLit.Elts {
				if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kvExpr.Key.(*ast.Ident); ok && ident.Name == fieldName {
						// Extract the value assigned to the field (e.g. "TestA" for `name: "TestA"`)
						nameValue := extractRHSValue(kvExpr.Value)
						nameValueStripped := strings.ReplaceAll(nameValue, `"`, "")
						cases = append(cases, &Case{
							Name: nameValueStripped,
						})
					}
				}
			}
		}
		return true
	})

	return cases
}
