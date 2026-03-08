package linter

import (
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func isError(t types.Type) bool {
	return types.Implements(t, errorType)
}

func isBlankIdent(expr ast.Expr) bool {
	ident, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && ident.Name == "_"
}

func exprToString(expr ast.Expr) string {
	var s strings.Builder
	err := printer.Fprint(&s, token.NewFileSet(), expr)
	if err != nil {
		panic(err)
	}
	return s.String()
}
