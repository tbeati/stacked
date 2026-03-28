package linter

import (
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"regexp"
	"strings"
)

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func implementsError(t types.Type) bool {
	return types.Implements(t, errorType)
}

func isError(t types.Type) bool {
	return t.String() == "error"
}

func isBool(t types.Type) bool {
	return types.Identical(t, types.Typ[types.Bool])
}

func isBlankIdent(expr ast.Expr) bool {
	ident, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && ident.Name == "_"
}

func isFunctionLiteral(call *ast.CallExpr) bool {
	_, isFuncLit := call.Fun.(*ast.FuncLit)
	return isFuncLit
}

func exprToString(expr ast.Expr) string {
	var s strings.Builder
	err := printer.Fprint(&s, token.NewFileSet(), expr)
	if err != nil {
		panic(err)
	}
	return s.String()
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
