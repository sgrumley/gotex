package finder

import (
	"fmt"
	"go/ast"
)

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
