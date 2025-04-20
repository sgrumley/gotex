package scanner

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strings"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/models"
	"github.com/sgrumley/gotex/pkg/slogger"
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
				cases, err := FindTestCases(ctx, fns[k])
				if err != nil {
					return nil, err
				}
				fns[k].Cases = cases
				for _, c := range cases {
					if c.Name != "" {
						fns[k].CaseMap[c.Name] = c
					}
				}
				fns[k].Parent = p.Packages[i].Files[j]
			}

			fns = mergeDuplicateFunctions(fns)
			p.Packages[i].Files[j].Functions = fns
			for _, fn := range fns {
				p.Packages[i].Files[j].FunctionMap[fn.Name] = fn
			}
		}
	}
	return p, nil
}

func FindTestFunctions(ctx context.Context, file *models.File) ([]*models.Function, error) {
	log, err := slogger.FromContext(ctx)
	if err != nil {
		log, _ = slogger.New(
			slogger.WithLevel(slog.LevelDebug),
			slogger.WithSource(false),
		)
	}

	ctx = slogger.AddToContext(ctx, log)
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
			log.Debug("test function found",
				slog.String("function name", fn.Name.Name),
				slog.String("file name", file.Name),
			)
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

// HACK: post processing the functions to combine any duplicates. This comes from multiple t.Run calls in the same test function
// This means that the VarName and RunCallNode elements of the Function are not accurate
// NOTE: this needs to be called after the cases have been found other wise it will only generate for one function
func mergeDuplicateFunctions(functions []*models.Function) []*models.Function {
	// Use map to track unique functions by name
	funcMap := make(map[string]*models.Function)

	// First pass: collect all functions and merge their cases
	for i := range functions {
		fn := functions[i]
		if existing, exists := funcMap[fn.Name]; exists {
			// Merge Cases
			existing.Cases = append(existing.Cases, fn.Cases...)

			// Merge CaseMap
			if existing.CaseMap == nil {
				existing.CaseMap = make(map[string]*models.Case)
			}
			for k, v := range fn.CaseMap {
				existing.CaseMap[k] = v
			}

			if existing.TestFunctionNode == nil {
				existing.TestFunctionNode = fn.TestFunctionNode
			}
		} else {
			// Create a new entry if it doesn't exist
			funcMap[fn.Name] = fn
		}
	}

	// Convert map back to slice
	result := make([]*models.Function, 0, len(funcMap))
	for _, fn := range funcMap {
		result = append(result, fn)
	}

	return result
}

func FindTestCases(ctx context.Context, fn *models.Function) ([]*models.Case, error) {
	log, err := slogger.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// exprToString(fn.RunCallNode.Args[0]) prints subtest e.g. name:tc.name // subtest name:tt.input // subtest name:name
	tRunParamType := ""
	if lit, ok := fn.RunCallNode.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		tRunParamType = "STRING"
	} else if _, ok := fn.RunCallNode.Args[0].(*ast.Ident); ok {
		tRunParamType = "IDENTIFIER"
	}

	tRunParam := exprToString(fn.RunCallNode.Args[0])
	tRunParamSplit := strings.Split(tRunParam, ".")
	var cases []*models.Case

	// this switches on if the variable for test name in t.Run(name, ...) contains a '.' inferring that it is an element of a struct
	switch len(tRunParamSplit) {
	case 1:
		// could be a static name ("success", BasicLit) or map key (name, Identifier)
		caseName := tRunParamSplit[0]
		switch tRunParamType {
		case "STRING":
			tc := &models.Case{
				Name:   caseName,
				Parent: fn,
				// Location: fn.RunCallNode.Lparen, // TODO: Location should be a more general type
			}
			fn.Cases = append(fn.Cases, tc)
			fn.CaseMap[tc.Name] = tc
			cases = append(cases, tc)
		case "IDENTIFIER":
			cases = append(cases, GetMapCase(log, fn, caseName)...)
		}
	case 2:
		// an element of a struct (tc.name)
		caseName := tRunParamSplit[1]
		cases = append(cases, GetElementCase(log, fn, caseName)...)
	default:
		log.Error("failed identifying first t.Run() param", nil,
			slog.String("closest param", tRunParam))
	}

	return cases, nil
}

func GetElementCase(log *slogger.Logger, fn *models.Function, caseName string) []*models.Case {
	var cases []*models.Case
	// looking for the value of a field in an array or slice (e.g. tc[i].name)
	ast.Inspect(fn.TestFunctionNode.Body, func(n ast.Node) bool {
		// We're looking for composite literals (array/slice initialization) or assignments
		if compLit, ok := n.(*ast.CompositeLit); ok {
			for _, elt := range compLit.Elts {
				if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kvExpr.Key.(*ast.Ident); ok && ident.Name == caseName {
						// Extract the value assigned to the field (e.g. "TestA" for `name: "TestA"`)
						nameValue := extractRHSValue(kvExpr.Value)
						nameValueStripped := strings.ReplaceAll(nameValue, `"`, "")
						log.Debug("test case found",
							slog.String("function name", fn.Name),
							slog.String("file name", fn.Parent.Name),
						)
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

func GetMapCase(log *slogger.Logger, fn *models.Function, caseName string) []*models.Case {
	var cases []*models.Case
	ast.Inspect(fn.TestFunctionNode.Body, func(n ast.Node) bool {
		// Looking for composite literals (map initialization)
		if compLit, ok := n.(*ast.CompositeLit); ok {
			for _, elt := range compLit.Elts {
				if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
					if bLit, ok := kvExpr.Key.(*ast.BasicLit); ok && bLit.Kind == token.STRING {
						nameValueStripped := strings.ReplaceAll(bLit.Value, `"`, "")
						log.Debug("test case found",
							slog.String("function name", fn.Name),
							slog.String("file name", fn.Parent.Name),
						)
						tc := &models.Case{
							Name:     nameValueStripped,
							Parent:   fn,
							Location: kvExpr,
						}
						fn.CaseMap[tc.Name] = tc
						fn.Cases = append(fn.Cases, tc)
						cases = append(cases, tc)
					}
					continue
				}
			}
			// since there is a nesting of key values on the structs as values
			// stop traversing once this level of the tree has been iterated
			return false
		}

		return true
	})

	return cases
}
