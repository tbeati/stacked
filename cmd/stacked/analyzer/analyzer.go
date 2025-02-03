package analyzer

import (
	"go/ast"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer(generatedPackages []string) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "stacked",
		Doc:  "check for error not wrapped with stacked",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			newChecker(generatedPackages, pass).check()
			return nil, nil
		},
	}
}

type checker struct {
	generatedPackages     []string
	pass                  *analysis.Pass
	endOfStatementCounter int
	toCheckIdent          *ast.Ident
}

func newChecker(generatedPackage []string, pass *analysis.Pass) *checker {
	return &checker{
		generatedPackages:     generatedPackage,
		pass:                  pass,
		endOfStatementCounter: -1,
	}
}

func (c *checker) check() {
	for _, file := range c.pass.Files {
		c.endOfStatementCounter = -1
		c.toCheckIdent = nil

		ast.Inspect(file, func(node ast.Node) bool {
			//fmt.Println("node", node, endOfStatementCounter)
			c.checkStmt(node)

			switch stmt := node.(type) {
			case *ast.AssignStmt:
				c.endOfStatementCounter = 1
				errors := c.handleAssignment(stmt)
				switch {
				case len(errors) > 1:
					c.pass.Reportf(stmt.Pos(), "multiple errors")
				case len(errors) == 1:
					c.toCheckIdent = errors[0]
				}
				//fmt.Println("errors", toCheckIdent)
			case *ast.IfStmt:
				//fmt.Println("if")
			case *ast.ReturnStmt:
				//fmt.Println("return")
			}

			return true
		})
	}
}

func (c *checker) checkStmt(node ast.Node) {
	if c.endOfStatementCounter == 0 {
		c.endOfStatementCounter = -1
	} else if c.endOfStatementCounter > -1 {
		if node == nil {
			c.endOfStatementCounter--
		} else {
			c.endOfStatementCounter++
		}
	}

	if c.toCheckIdent != nil && c.endOfStatementCounter == -1 {
		toCheckIdent := c.toCheckIdent
		c.toCheckIdent = nil

		assignStmt, ok := node.(*ast.AssignStmt)
		if !ok {
			c.report(toCheckIdent)
			return
		}

		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			c.report(toCheckIdent)
			return
		}

		call, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			c.report(toCheckIdent)
			return
		}

		if !c.isStackTraceWrap(call) {
			c.report(toCheckIdent)
			return
		}

	}
}

func (c *checker) report(ident *ast.Ident) {
	c.pass.Reportf(ident.Pos(), "%s is not wrapped with stacked", ident.Name)
}

func (c *checker) handleAssignment(stmt *ast.AssignStmt) []*ast.Ident {
	var errorIdents []*ast.Ident

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

			if c.isInternalCall(call) {
				continue
			}

			if c.isStackTraceWrap(call) {
				continue
			}

			if c.returnsError(call) {
				errorIdents = append(errorIdents, ident)
			}
		}
	} else {
		call, ok := stmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return nil
		}

		if c.isInternalCall(call) {
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
				errorIdents = append(errorIdents, ident)
			}
		}
	}

	return errorIdents
}

func isError(t types.Type) bool {
	return t.String() == "error"
}

func (c *checker) isStackTraceWrap(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if selector.Sel.Name != "Wrap" {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := c.pass.TypesInfo.Uses[ident]

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkg.Imported().Path() == "github.com/beati/stacked"
}

func (c *checker) isInternalCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	pkg := c.pass.TypesInfo.Uses[selector.Sel].Pkg().Path()

	if slices.Contains(c.generatedPackages, pkg) {
		return false
	}

	return strings.HasPrefix(pkg, c.pass.Module.Path)
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
