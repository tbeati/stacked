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

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

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
	generatedPackages []string
	pass              *analysis.Pass

	endOfFunctionCounter int

	endOfStatementCounter int
	assignedError         ast.Expr
	toCheckExpr           ast.Expr
}

func newChecker(generatedPackage []string, pass *analysis.Pass) *checker {
	return &checker{
		generatedPackages:     generatedPackage,
		pass:                  pass,
		endOfStatementCounter: -1,
	}
}

func (c *checker) reset() {
	c.endOfFunctionCounter = -1
	c.endOfStatementCounter = -1
	c.assignedError = nil
	c.toCheckExpr = nil
}

func (c *checker) isInFunction() bool {
	return c.endOfFunctionCounter > -1
}

func (c *checker) updateFunctionState(node ast.Node) {
	if c.endOfFunctionCounter == 0 {
		c.endOfFunctionCounter = -1
	} else if c.endOfFunctionCounter > -1 {
		if node == nil {
			c.endOfFunctionCounter--
		} else {
			c.endOfFunctionCounter++
		}
	}

	if c.isInFunction() {
		return
	}

	switch node.(type) {
	case *ast.FuncDecl:
		c.endOfFunctionCounter = 1
	case *ast.FuncLit:
		c.endOfFunctionCounter = 1
	}
}

func (c *checker) check() {
	for _, file := range c.pass.Files {
		c.reset()

		ast.Inspect(file, func(node ast.Node) bool {
			c.updateFunctionState(node)
			c.checkStmt(node)

			if !c.isInFunction() {
				return true
			}

			switch node := node.(type) {
			case *ast.AssignStmt:
				c.checkAssignStmt(node)
			case *ast.ReturnStmt:
				for _, result := range node.Results {
					c.reportShouldWrap(result)
				}
			case *ast.CompositeLit:
				for i := range node.Elts {
					switch elt := ast.Unparen(node.Elts[i]).(type) {
					case *ast.KeyValueExpr:
						c.reportShouldWrap(elt.Value)
					default:

						c.reportShouldWrap(elt)
					}
				}
			case *ast.CallExpr:
				if c.isStackedWrap(node) {
					break
				}

				for _, arg := range node.Args {
					c.reportShouldWrap(arg)
				}
			}

			return true
		})
	}
}

func (c *checker) shouldWrap(expr ast.Expr) bool {
	expr = ast.Unparen(expr)

	switch expr := expr.(type) {
	case *ast.Ident:
		return c.shouldWrapIdent(expr)
	case *ast.SelectorExpr:
		return c.shouldWrapSelector(expr)
	case *ast.CompositeLit:
		return c.shouldWrapCompositeLit(expr)
	case *ast.CallExpr:
		return c.shouldWrapCall(expr)
	}

	return false
}

func (c *checker) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return isError(c.pass.TypesInfo.TypeOf(lit))
}

func (c *checker) shouldWrapIdent(ident *ast.Ident) bool {
	obj := c.pass.TypesInfo.Uses[ident]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (c *checker) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	obj := c.pass.TypesInfo.Uses[expr.Sel]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (c *checker) shouldWrapCall(call *ast.CallExpr) bool {
	if c.isInternalCall(call) {
		return false
	}

	if c.isStackedWrap(call) {
		return false
	}

	return c.returnsError(call)
}

func (c *checker) reportShouldWrap(expr ast.Expr) {
	if c.shouldWrap(expr) {
		c.report(expr)
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

	if c.assignedError != nil && c.endOfStatementCounter == -1 {
		assignedError := c.assignedError
		toCheckCall := c.toCheckExpr
		c.assignedError = nil
		c.toCheckExpr = nil

		assignStmt, ok := node.(*ast.AssignStmt)
		if !ok {
			c.report(toCheckCall)
			return
		}

		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			c.report(toCheckCall)
			return
		}

		call, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			c.report(toCheckCall)
			return
		}

		if !c.isStackedWrap(call) {
			c.report(toCheckCall)
			return
		}

		if !areExprsEqual(assignStmt.Lhs[0], assignedError) {
			c.report(toCheckCall)
			return
		}

		if len(call.Args) != 1 {
			c.report(toCheckCall)
			return
		}

		if !areExprsEqual(call.Args[0], assignedError) {
			c.report(toCheckCall)
			return
		}
	}
}

func (c *checker) report(expr ast.Expr) {
	switch expr := expr.(type) {
	case *ast.Ident:
		c.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.SelectorExpr:
		c.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.CompositeLit:
		c.pass.Reportf(expr.Pos(), "%s literal is not wrapped with stacked", exprToString(expr.Type))
	case *ast.CallExpr:
		c.pass.Reportf(expr.Pos(), "error returned by %s is not wrapped with stacked", exprToString(expr.Fun))
	}
}

func (c *checker) checkAssignStmt(stmt *ast.AssignStmt) {
	c.endOfStatementCounter = 1

	errCount := 0
	for _, expr := range stmt.Lhs {
		exprType := c.pass.TypesInfo.TypeOf(expr)
		if exprType != nil && isError(exprType) {
			errCount++
		}

		if errCount > 1 {
			c.pass.Reportf(stmt.Pos(), "multiple errors")
			return
		}
	}

	if len(stmt.Lhs) == len(stmt.Rhs) {
		for i := range stmt.Lhs {
			if c.shouldWrap(stmt.Rhs[i]) {
				c.assignedError = ast.Unparen(stmt.Lhs[i])
				c.toCheckExpr = ast.Unparen(stmt.Rhs[i])
				return
			}
		}
	} else {
		call, ok := ast.Unparen(stmt.Rhs[0]).(*ast.CallExpr)
		if !ok {
			return
		}

		if c.isInternalCall(call) {
			return
		}

		if c.isStackedWrap(call) {
			return
		}

		errors := c.errorsByArg(call)
		for i := range stmt.Lhs {
			if errors[i] {
				c.assignedError = ast.Unparen(stmt.Lhs[i])
				c.toCheckExpr = call
				return
			}
		}
	}
}

func isError(t types.Type) bool {
	return types.Implements(t, errorType)
}

func (c *checker) isStackedWrap(call *ast.CallExpr) bool {
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
