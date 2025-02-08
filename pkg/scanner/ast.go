package scanner

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

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
func findEnclosingFunction(node ast.Node, n ast.Node) *ast.FuncDecl {
	// Walk the AST to find the wrapping function
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

// HELPERS

// exprToString converts an AST expression into a string
func exprToString(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, "\"")
	}
	return formatExpr(expr) // Dynamic name
}

// extractRHSValue returns a string representation of the right-hand side of an assignment
func extractRHSValue(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		// Handle basic literals (e.g., string or number)
		return v.Value
	case *ast.Ident:
		// Handle identifiers (e.g., variables)
		return v.Name
	case *ast.BinaryExpr:
		// Handle binary expressions (e.g., concatenation like "prefix" + "suffix")
		return fmt.Sprintf("%s %s %s", extractRHSValue(v.X), v.Op, extractRHSValue(v.Y))
	default:
		return fmt.Sprintf("unknown value (%T)", expr)
	}
}

// formatExpr returns a string representation of an expression in the AST
func formatExpr(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name // Variable name
	case *ast.BinaryExpr:
		// Handle binary expressions (e.g., concatenation like var1 + var2)
		return fmt.Sprintf("(%s %s %s)", formatExpr(v.X), v.Op.String(), formatExpr(v.Y))
	case *ast.BasicLit:
		return v.Value // Literal value
	case *ast.CallExpr:
		// Handle function calls (e.g., fmt.Sprintf)
		return fmt.Sprintf("function call: %s", v.Fun)
	case *ast.SelectorExpr:
		// Handle selector expressions (e.g., pkg.var)
		return fmt.Sprintf("%s.%s", formatExpr(v.X), v.Sel.Name)
	default:
		return fmt.Sprintf("unknown expression type (%T)", expr)
	}
}
