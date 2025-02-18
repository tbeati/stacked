package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	GeneratedPackages []string
}

func (c *Config) isGeneratedPackage(pkg string) bool {
	for _, genPkg := range c.GeneratedPackages {
		if strings.HasPrefix(pkg, genPkg) {
			return true
		}
	}

	return false
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "stacked",
		Doc:  "check for error not wrapped with stacked",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			if config.isGeneratedPackage(pass.Pkg.Path()) {
				return nil, nil
			}

			for _, file := range pass.Files {
				newFileChecker(config, pass, file).check()
			}
			return nil, nil
		},
	}
}

type fileChecker struct {
	config *Config
	pass   *analysis.Pass
	file   *ast.File

	functionTracker   nodeTracker
	assignmentTracker nodeTracker
	assignedErrorDst  ast.Expr
	assignedErrorSrc  ast.Expr
}

func newFileChecker(config *Config, pass *analysis.Pass, file *ast.File) *fileChecker {
	return &fileChecker{
		config:            config,
		pass:              pass,
		file:              file,
		functionTracker:   newNodeTracker(),
		assignmentTracker: newNodeTracker(),
	}
}

func (c *fileChecker) check() {
	ast.Inspect(c.file, func(node ast.Node) bool {
		c.functionTracker.step(node)
		c.assignmentTracker.step(node)
		c.trackTopLevelFunctionDeclaration(node)

		c.checkAssignmentWrapping(node)

		if !c.isInFunction() {
			return true
		}

		switch node := node.(type) {
		case *ast.DeclStmt:
			c.assignmentTracker.enter()
		case *ast.GenDecl:
			c.checkGenDecl(node)
		case *ast.AssignStmt:
			c.assignmentTracker.enter()
			c.checkAssignStmt(node)
		case *ast.ReturnStmt:
			for _, result := range node.Results {
				if c.shouldWrap(result) {
					c.report(result)
				}
			}
		case *ast.CompositeLit:
			for _, elt := range node.Elts {
				switch elt := ast.Unparen(elt).(type) {
				case *ast.KeyValueExpr:
					if c.shouldWrap(elt.Value) {
						c.report(elt.Value)
					}
				default:
					if c.shouldWrap(elt) {
						c.report(elt)
					}
				}
			}
		case *ast.CallExpr:
			if c.isWrapCall(node) {
				return true
			}

			for _, arg := range node.Args {
				if c.shouldWrap(arg) {
					c.report(arg)
				}
			}
		}

		return true
	})
}

func (c *fileChecker) isInFunction() bool {
	return c.functionTracker.isIn()
}

func (c *fileChecker) trackTopLevelFunctionDeclaration(node ast.Node) {
	if c.isInFunction() {
		return
	}

	switch node.(type) {
	case *ast.FuncDecl:
		c.functionTracker.enter()
	case *ast.FuncLit:
		c.functionTracker.enter()
	}
}

func (c *fileChecker) report(expr ast.Expr) {
	switch expr := expr.(type) {
	case *ast.Ident:
		c.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.SelectorExpr:
		c.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.CompositeLit:
		c.pass.Reportf(expr.Pos(), "%s literal is not wrapped with stacked", exprToString(expr.Type))
	case *ast.CallExpr:
		if c.isTypeConversion(expr) {
			c.pass.Reportf(expr.Pos(), "value converted to error type %s is not wrapped with stacked", exprToString(expr.Fun))
		} else {
			c.pass.Reportf(expr.Pos(), "error returned by %s is not wrapped with stacked", exprToString(expr.Fun))
		}
	}
}

func (c *fileChecker) shouldWrap(expr ast.Expr) bool {
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

func (c *fileChecker) shouldWrapIdent(ident *ast.Ident) bool {
	obj := c.pass.TypesInfo.Uses[ident]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (c *fileChecker) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	obj := c.pass.TypesInfo.Uses[expr.Sel]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (c *fileChecker) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return isError(c.pass.TypesInfo.TypeOf(lit))
}

func (c *fileChecker) shouldWrapCall(call *ast.CallExpr) bool {
	if c.isInternalCall(call) {
		return false
	}

	if c.isWrapCall(call) {
		return false
	}

	return c.returnsError(call)
}

func (c *fileChecker) checkGenDecl(stmt *ast.GenDecl) {
	if stmt.Tok != token.VAR {
		return
	}

	var errorSpec *ast.ValueSpec

	errCount := 0
	for _, spec := range stmt.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, name := range valueSpec.Names {
			if isError(c.pass.TypesInfo.TypeOf(name)) {
				errCount++
				if errCount > 1 {
					c.pass.Reportf(stmt.Pos(), "multiple errors")
					return
				}

				errorSpec = valueSpec
			}
		}
	}

	if errorSpec == nil || len(errorSpec.Values) == 0 {
		return
	}

	lsh := make([]ast.Expr, 0, len(errorSpec.Names))
	for _, ident := range errorSpec.Names {
		lsh = append(lsh, ident)
	}

	c.checkAssignment(lsh, errorSpec.Values)
}

func (c *fileChecker) checkAssignStmt(stmt *ast.AssignStmt) {
	errCount := 0
	for _, expr := range stmt.Lhs {
		exprType := c.pass.TypesInfo.TypeOf(expr)
		if exprType != nil && isError(exprType) {
			errCount++
			if errCount > 1 {
				c.pass.Reportf(stmt.Pos(), "multiple errors")
				return
			}
		}
	}

	c.checkAssignment(stmt.Lhs, stmt.Rhs)
}

func (c *fileChecker) checkAssignment(lsh, rsh []ast.Expr) {
	if len(lsh) == len(rsh) {
		for i := range lsh {
			if c.shouldWrap(rsh[i]) {
				c.assignedErrorDst = ast.Unparen(lsh[i])
				c.assignedErrorSrc = ast.Unparen(rsh[i])
				return
			}
		}
	} else {
		call, ok := ast.Unparen(rsh[0]).(*ast.CallExpr)
		if !ok {
			return
		}

		if !c.shouldWrapCall(call) {
			return
		}

		c.assignedErrorDst = ast.Unparen(lsh[c.errorReturnIndex(call)])
		c.assignedErrorSrc = call
	}
}

func (c *fileChecker) checkAssignmentWrapping(node ast.Node) {
	if c.assignedErrorDst != nil && !c.assignmentTracker.isIn() {
		assignedErrorDst := c.assignedErrorDst
		assignedErrorSrc := c.assignedErrorSrc
		c.assignedErrorDst = nil
		c.assignedErrorSrc = nil

		assignStmt, ok := node.(*ast.AssignStmt)
		if !ok {
			c.report(assignedErrorSrc)
			return
		}

		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			c.report(assignedErrorSrc)
			return
		}

		call, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			c.report(assignedErrorSrc)
			return
		}

		if !c.isWrapCall(call) {
			c.report(assignedErrorSrc)
			return
		}

		if !areExprsEqual(assignStmt.Lhs[0], assignedErrorDst) {
			c.report(assignedErrorSrc)
			return
		}

		if len(call.Args) != 1 {
			c.report(assignedErrorSrc)
			return
		}

		if !areExprsEqual(call.Args[0], assignedErrorDst) {
			c.report(assignedErrorSrc)
			return
		}
	}
}

func (c *fileChecker) isWrapCall(call *ast.CallExpr) bool {
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

func (c *fileChecker) isInternalCall(call *ast.CallExpr) bool {
	if c.isTypeConversion(call) {
		return false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	pkg := c.pass.TypesInfo.Uses[selector.Sel].Pkg()
	if pkg == nil {
		return false
	}

	if c.config.isGeneratedPackage(pkg.Path()) {
		return false
	}

	return strings.HasPrefix(pkg.Path(), c.pass.Module.Path)
}

func (c *fileChecker) isTypeConversion(call *ast.CallExpr) bool {
	var obj types.Object

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		obj = c.pass.TypesInfo.Uses[fun]
	case *ast.SelectorExpr:
		obj = c.pass.TypesInfo.Uses[fun.Sel]
	}

	if obj == nil {
		return false
	}

	_, isTypeName := obj.(*types.TypeName)
	return isTypeName
}

func (c *fileChecker) errorReturnIndex(call *ast.CallExpr) int {
	switch returnType := c.pass.TypesInfo.TypeOf(call).(type) {
	case *types.Named:
		if isError(returnType) {
			return 0
		}
	case *types.Tuple:
		for i := range returnType.Len() {
			t, ok := returnType.At(i).Type().(*types.Named)
			if ok && isError(t) {
				return i
			}
		}
	}

	return -1
}

func (c *fileChecker) returnsError(call *ast.CallExpr) bool {
	return c.errorReturnIndex(call) >= 0
}
