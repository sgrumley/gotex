package finder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strings"
)

func SearchFile(file *File, log *slog.Logger) error {
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, file.Path, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", file.Path, err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		// Check all call expressions (function calls)
		if callExpr, ok := n.(*ast.CallExpr); ok {
			// a selector expression is the syntax used to access fields or methods of a struct
			// checking that the selector is `Run`, we can then check that the expression is `t`
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "Run" {
				if ex := exprToString(selExpr.X); ex != "t" {
					log.Debug("here", slog.String("expression", ex))
					return true
				}
				// Get the parent function (this is intended to be the test function containing the subtests)
				fn := findEnclosingFunction(node, callExpr)
				if fn == nil {
					return true
				}

				// Extract the subtest variable name (this should be something like tc.name where tc is the stuct and name is the attribute provided to t.Run(tc.name, ...))
				subtestName := exprToString(callExpr.Args[0])
				log.Debug("test case found",
					slog.String("case name", subtestName),
					slog.String("function name", fn.Name.Name),
					slog.String("file name", file.Name),
				)

				// setting a default that I use, probably needs better error handling in general
				// structName := "tc"
				caseName := "name"

				subtestNameSplit := strings.Split(subtestName, ".")
				if len(subtestNameSplit) == 2 {
					// structName = subtestNameSplit[0]
					caseName = subtestNameSplit[1]
				} else {
					log.Error("failed identifying struct.name, defaulting to tc.name", slog.String("name", subtestName))
				}

				// TODO: add the name of the struct back in for validation
				// tcNames := findValuesOfIndexedField(fn, "tc", "name")

				// Find all occurrences of `tc.name` in the function
				cases := findValuesOfIndexedField(fn, caseName)
				caseMap := make(map[string]*Case)
				f := &Function{}

				for i := range cases {
					cases[i].Parent = f
					caseMap[cases[i].Name] = cases[i]
				}

				f.Name = fn.Name.Name
				f.Cases = cases
				f.CaseMap = caseMap
				f.decl = fn
				f.VarName = subtestName
				f.Parent = file

				file.FunctionMap[fn.Name.Name] = f
				file.Functions = append(file.Functions, f)
			}
		}
		return true
	})

	return nil
}

// extractSubtestName handles both string literals and dynamic subtest names
func exprToString(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, "\"")
	}
	return formatExpr(expr) // Dynamic name
}

// findEnclosingFunction traverses the AST upwards to find the function that encloses the node
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
