package scanner

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/models"
)

func Scan(ctx context.Context, cfg config.Config, root string) (*models.Project, error) {
	p := &models.Project{
		Config: cfg,
		// Log:      log,
		RootDir: root,
	}

	pkgs, err := FindPackagesAndTestFiles(root, p)
	if err != nil {
		return nil, err
	}
	p.Packages = pkgs

	for i := range p.Packages {
		// TODO: concurrent
		for j := range p.Packages[i].Files {
			p.Packages[i].Files[j].FunctionMap = make(map[string]*models.Function)

			fns, err := FindTestFunctions(ctx, p.Packages[i].Files[j])
			if err != nil {
				return nil, fmt.Errorf("failed finding test functions in file: %s: %w", p.Packages[i].Files[j].Path, err)
			}

			for k := range fns {
				fns[k].CaseMap = make(map[string]*models.Case)
				cases := FindTestCases(ctx, fns[k])
				fns[k].Cases = cases
				for _, c := range cases {
					if c.Name != "" {
						fns[k].CaseMap[c.Name] = c
					}
				}
				fns[k].Parent = p.Packages[i].Files[j]
			}

			p.Packages[i].Files[j].Functions = fns
			for _, fn := range fns {
				p.Packages[i].Files[j].FunctionMap[fn.Name] = fn
			}
		}
	}
	return p, nil
}

func FindTestFunctions(ctx context.Context, file *models.File) ([]*models.Function, error) {
	fset := token.NewFileSet()
	fns := make([]*models.Function, 0)
	node, err := parser.ParseFile(fset, file.Path, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", file.Path, err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			_, ok := isTestRunCall(callExpr)
			if !ok {
				return false
			}
			fn := findEnclosingFunction(node, callExpr)

			if fn == nil {
				return false
			}
			// from ctx
			// log.Debug("test case found",
			// 		slog.String("case name", subtestName),
			// 		slog.String("function name", fn.Name.Name),
			// 		slog.String("file name", file.Name),
			// )
			newFn := &models.Function{
				Name:             fn.Name.Name,
				TestFunctionNode: fn,
				RunCallNode:      callExpr,
				Parent:           file,
				CaseMap:          make(map[string]*models.Case),
			}

			file.FunctionMap[newFn.Name] = newFn
			file.Functions = append(file.Functions, newFn)
			fns = append(fns, newFn)
		}
		return true
	})

	return fns, nil
}

func FindTestCases(ctx context.Context, fn *models.Function) []*models.Case {
	caseName := extractCaseName(ctx, exprToString(fn.RunCallNode.Args[0]))
	var cases []*models.Case

	// findValuesOfIndexedField looks for the value of a field in an array or slice (e.g. tc[i].name)
	ast.Inspect(fn.TestFunctionNode.Body, func(n ast.Node) bool {
		// We're looking for composite literals (array/slice initialization) or assignments
		if compLit, ok := n.(*ast.CompositeLit); ok {
			for _, elt := range compLit.Elts {
				if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kvExpr.Key.(*ast.Ident); ok && ident.Name == caseName {
						// Extract the value assigned to the field (e.g. "TestA" for `name: "TestA"`)
						nameValue := extractRHSValue(kvExpr.Value)
						nameValueStripped := strings.ReplaceAll(nameValue, `"`, "")
						tc := &models.Case{
							Name:     nameValueStripped,
							Parent:   fn,
							Location: kvExpr,
						}
						fn.CaseMap[tc.Name] = tc
						fn.Cases = append(fn.Cases, tc)
						cases = append(cases, tc)

					}
				}
			}
		}
		return true
	})

	return cases
}

// extractCaseName gets the case field name from the subtest name
func extractCaseName(ctx context.Context, subtestName string) string {
	caseName := "name"
	subtestNameSplit := strings.Split(subtestName, ".")

	if len(subtestNameSplit) == 2 {
		caseName = subtestNameSplit[1]
	} else {
		// from ctx
		// log.Error("failed identifying struct.name, defaulting to tc.name",
		// 	slog.String("name", subtestName))
	}

	return caseName
}
