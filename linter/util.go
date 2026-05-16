package linter

import (
	"go/ast"
	"go/types"
	"regexp"
	"strings"
)

var errorType types.Type
var errorTypeInterface *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type()
	errorTypeInterface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func implementsError(t types.Type) bool {
	return types.Implements(t, errorTypeInterface)
}

func isError(t types.Type) bool {
	return types.Identical(t, errorType)
}

func isBool(t types.Type) bool {
	return types.Identical(t, types.Typ[types.Bool])
}

func isBlankIdent(expr ast.Expr) bool {
	ident, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && ident.Name == "_"
}

func typeToString(t types.Type, currentPkg *types.Package) string {
	qualifier := func(pkg *types.Package) string {
		if pkg == currentPkg {
			return ""
		}
		return pkg.Name()
	}

	return types.TypeString(t, qualifier)
}

var wVerbRegex = regexp.MustCompile(`%[^a-zA-Z%]*w`)

func containsWVerb(formatString string) bool {
	cleanFormat := strings.ReplaceAll(formatString, "%%", "")
	return wVerbRegex.MatchString(cleanFormat)
}

func stripTypeArgs(expr ast.Expr) ast.Expr {
	expr = ast.Unparen(expr)
	switch x := expr.(type) {
	case *ast.IndexExpr:
		return ast.Unparen(x.X)
	case *ast.IndexListExpr:
		return ast.Unparen(x.X)
	}
	return expr
}
