package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "stacked",
	Doc:  "check for error not wrapped with stacked",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		newChecker(pass).check()
		return nil, nil
	},
}

type checker struct {
	pass *analysis.Pass
}

func newChecker(pass *analysis.Pass) *checker {
	return &checker{
		pass: pass,
	}
}

func (c *checker) check() {
	for _, file := range c.pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.AssignStmt:
				fmt.Println(stmt)
				idents := c.handleAssignment(stmt)
				fmt.Println(idents)
			}

			return true
		})
	}
}

func (c *checker) handleAssignment(stmt *ast.AssignStmt) []string {
	var errorIdents []string

	if len(stmt.Lhs) == len(stmt.Rhs) {
		for i := range stmt.Lhs {
			ident, ok := stmt.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}

			call, ok := stmt.Rhs[i].(*ast.CallExpr)
			if !ok {
				continue
			}

			if c.isStackTraceWrap(call) {
				continue
			}

			if c.returnsError(call) {
				errorIdents = append(errorIdents, ident.Name)
			}
		}
	} else {
		call, ok := stmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return nil
		}

		if c.isStackTraceWrap(call) {
			return nil
		}

		errors := c.errorsByArg(call)
		for i := range stmt.Lhs {
			ident, ok := stmt.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}

			if errors[i] {
				errorIdents = append(errorIdents, ident.Name)
			}
		}
	}

	return errorIdents
}

func isError(t types.Type) bool {
	return t.String() == "error"
}

func (c *checker) isStackTraceWrap(call *ast.CallExpr) bool {
	fmt.Println("func", call.Fun)
	return false
}

func (c *checker) errorsByArg(call *ast.CallExpr) []bool {
	switch t := c.pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named:
		return []bool{isError(t)}
	case *types.Tuple:
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch et := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isError(et)
			default:
				s[i] = false
			}
		}

		return s
	}

	return []bool{false}
}

func (c *checker) returnsError(call *ast.CallExpr) bool {
	for _, isError := range c.errorsByArg(call) {
		if isError {
			return true
		}
	}
	return false
}
