package analyzer

import (
	"go/ast"
	"go/printer"
	"go/token"
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
	toCheckIdent          ast.Expr
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
				c.toCheckIdent = c.handleAssignment(stmt)
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

		if !areAssignableExprEqual(assignStmt.Lhs[0], toCheckIdent) {
			c.report(toCheckIdent)
			return
		}

		if len(call.Args) != 1 {
			c.report(toCheckIdent)
			return
		}

		if !areAssignableExprEqual(call.Args[0], toCheckIdent) {
			c.report(toCheckIdent)
			return
		}
	}
}

func (c *checker) report(errorExpr ast.Expr) {
	c.pass.Reportf(errorExpr.Pos(), "%s is not wrapped with stacked", exprToString(errorExpr))
}

func (c *checker) handleAssignment(stmt *ast.AssignStmt) ast.Expr {
	var assignedErrorExpr ast.Expr

	if len(stmt.Lhs) == len(stmt.Rhs) {
		for i := range stmt.Lhs {
			call, ok := ast.Unparen(stmt.Rhs[i]).(*ast.CallExpr)
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
				if assignedErrorExpr != nil {
					c.pass.Reportf(stmt.Pos(), "multiple errors")
					return nil
				}

				assignedErrorExpr = ast.Unparen(stmt.Lhs[i])
			}
		}
	} else {
		call, ok := ast.Unparen(stmt.Rhs[0]).(*ast.CallExpr)
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
			if errors[i] {
				if assignedErrorExpr != nil {
					c.pass.Reportf(stmt.Pos(), "multiple errors")
					return nil
				}

				assignedErrorExpr = ast.Unparen(stmt.Lhs[i])
			}
		}
	}

	return assignedErrorExpr
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
	return slices.Contains(c.errorsByArg(call), true)
}

func exprToString(expr ast.Expr) string {
	var sb strings.Builder
	err := printer.Fprint(&sb, token.NewFileSet(), expr)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

func areAssignableExprEqual(a, b ast.Expr) bool {
	switch a := a.(type) {
	case *ast.Ident:
		b, ok := b.(*ast.Ident)
		return ok && a.Name == b.Name
	case *ast.BasicLit:
		b, ok := b.(*ast.BasicLit)
		return ok && a.Kind == b.Kind && a.Value == b.Value
	case *ast.CompositeLit:
		b, ok := b.(*ast.CompositeLit)
		if !ok || len(a.Elts) != len(b.Elts) {
			return false
		}
		if !areAssignableExprEqual(a.Type, b.Type) {
			return false
		}
		for i := range a.Elts {
			if !areAssignableExprEqual(a.Elts[i], b.Elts[i]) {
				return false
			}
		}
		return false
	case *ast.UnaryExpr:
		b, ok := b.(*ast.UnaryExpr)
		return ok && a.Op == b.Op && areAssignableExprEqual(a.X, b.X)
	case *ast.BinaryExpr:
		b, ok := b.(*ast.BinaryExpr)
		return ok && a.Op == b.Op && areAssignableExprEqual(a.X, b.X) && areAssignableExprEqual(a.Y, b.Y)
	case *ast.CallExpr:
		b, ok := b.(*ast.CallExpr)
		if !ok || len(a.Args) != len(b.Args) {
			return false
		}
		if !areAssignableExprEqual(a.Fun, b.Fun) {
			return false
		}
		for i := range a.Args {
			if !areAssignableExprEqual(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	case *ast.ParenExpr:
		b, ok := b.(*ast.ParenExpr)
		return ok && areAssignableExprEqual(a.X, b.X)
	case *ast.SelectorExpr:
		b, ok := b.(*ast.SelectorExpr)
		return ok && a.Sel.Name == b.Sel.Name && areAssignableExprEqual(a.X, b.X)
	case *ast.IndexExpr:
		b, ok := b.(*ast.IndexExpr)
		return ok && areAssignableExprEqual(a.X, b.X) && areAssignableExprEqual(a.Index, b.Index)
	case *ast.IndexListExpr:
		b, ok := b.(*ast.IndexListExpr)
		if !ok || len(a.Indices) != len(b.Indices) {
			return false
		}
		if !areAssignableExprEqual(a.X, b.X) {
			return false
		}
		for i := range a.Indices {
			if !areAssignableExprEqual(a.Indices[i], b.Indices[i]) {
				return false
			}
		}
		return true
	case *ast.SliceExpr:
		b, ok := b.(*ast.SliceExpr)
		return ok && a.Slice3 == b.Slice3 &&
			areAssignableExprEqual(a.X, b.X) &&
			areAssignableExprEqual(a.Low, b.Low) &&
			areAssignableExprEqual(a.High, b.High) &&
			areAssignableExprEqual(a.Max, b.Max)
	case *ast.KeyValueExpr:
		b, ok := b.(*ast.KeyValueExpr)
		return ok && areAssignableExprEqual(a.Key, b.Key) && areAssignableExprEqual(a.Value, b.Value)
	case *ast.TypeAssertExpr:
		b, ok := b.(*ast.TypeAssertExpr)
		return ok && areAssignableExprEqual(a.X, b.X) && areAssignableExprEqual(a.Type, b.Type)
	case *ast.StarExpr:
		b, ok := b.(*ast.StarExpr)
		return ok && areAssignableExprEqual(a.X, b.X)
	}

	return false
}
