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

func (fc *fileChecker) check() {
	ast.Inspect(fc.file, func(node ast.Node) bool {
		fc.functionTracker.depthFirstSearchStep(node)
		fc.assignmentTracker.depthFirstSearchStep(node)
		fc.trackTopLevelFunctionDeclaration(node)

		fc.checkAssignmentWrapping(node)

		if !fc.isInFunction() {
			return true
		}

		switch node := node.(type) {
		case *ast.DeclStmt:
			fc.assignmentTracker.enterNode()
		case *ast.GenDecl:
			fc.checkGenDecl(node)
		case *ast.AssignStmt:
			fc.assignmentTracker.enterNode()
			fc.checkAssignStmt(node)
		case *ast.ReturnStmt:
			for _, result := range node.Results {
				if fc.shouldWrap(result) {
					fc.report(result)
				}
			}
		case *ast.CompositeLit:
			for _, elt := range node.Elts {
				switch elt := ast.Unparen(elt).(type) {
				case *ast.KeyValueExpr:
					if fc.shouldWrap(elt.Value) {
						fc.report(elt.Value)
					}
				default:
					if fc.shouldWrap(elt) {
						fc.report(elt)
					}
				}
			}
		case *ast.CallExpr:
			if fc.isWrapCall(node) {
				return true
			}

			for _, arg := range node.Args {
				if fc.shouldWrap(arg) {
					fc.report(arg)
				}
			}
		}

		return true
	})
}

func (fc *fileChecker) isInFunction() bool {
	return fc.functionTracker.isInNode()
}

func (fc *fileChecker) trackTopLevelFunctionDeclaration(node ast.Node) {
	if fc.isInFunction() {
		return
	}

	switch node.(type) {
	case *ast.FuncDecl:
		fc.functionTracker.enterNode()
	case *ast.FuncLit:
		fc.functionTracker.enterNode()
	}
}

func (fc *fileChecker) report(expr ast.Expr) {
	switch expr := expr.(type) {
	case *ast.Ident:
		fc.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.SelectorExpr:
		fc.pass.Reportf(expr.Pos(), "%s is not wrapped with stacked", exprToString(expr))
	case *ast.CompositeLit:
		fc.pass.Reportf(expr.Pos(), "%s literal is not wrapped with stacked", exprToString(expr.Type))
	case *ast.CallExpr:
		if fc.isTypeConversion(expr) {
			fc.pass.Reportf(expr.Pos(), "value converted to error type %s is not wrapped with stacked", exprToString(expr.Fun))
		} else {
			fc.pass.Reportf(expr.Pos(), "error returned by %s is not wrapped with stacked", exprToString(expr.Fun))
		}
	}
}

func (fc *fileChecker) shouldWrap(expr ast.Expr) bool {
	if fc.isIgnored(expr) {
		return false
	}

	expr = ast.Unparen(expr)

	switch expr := expr.(type) {
	case *ast.Ident:
		return fc.shouldWrapIdent(expr)
	case *ast.SelectorExpr:
		return fc.shouldWrapSelector(expr)
	case *ast.CompositeLit:
		return fc.shouldWrapCompositeLit(expr)
	case *ast.CallExpr:
		return fc.shouldWrapCall(expr)
	}

	return false
}

func (fc *fileChecker) isIgnored(expr ast.Expr) bool {
	line := fc.pass.Fset.Position(expr.Pos()).Line

	for _, commentGroup := range fc.file.Comments {
		for _, comment := range commentGroup.List {
			commentLine := fc.pass.Fset.Position(comment.Pos()).Line
			if line == commentLine && strings.Contains(comment.Text, "//stacked:disable") {
				return true
			}
		}
	}

	return false
}

func (fc *fileChecker) shouldWrapIdent(ident *ast.Ident) bool {
	obj := fc.pass.TypesInfo.Uses[ident]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (fc *fileChecker) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	obj := fc.pass.TypesInfo.Uses[expr.Sel]

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope() && isError(variable.Type())
}

func (fc *fileChecker) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return isError(fc.pass.TypesInfo.TypeOf(lit))
}

func (fc *fileChecker) shouldWrapCall(call *ast.CallExpr) bool {
	if fc.isInternalCall(call) {
		return false
	}

	if fc.isWrapCall(call) {
		return false
	}

	return fc.returnsError(call)
}

func (fc *fileChecker) checkGenDecl(stmt *ast.GenDecl) {
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
			if isError(fc.pass.TypesInfo.TypeOf(name)) {
				errCount++
				if errCount > 1 {
					fc.pass.Reportf(stmt.Pos(), "multiple errors")
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

	fc.checkAssignment(lsh, errorSpec.Values)
}

func (fc *fileChecker) checkAssignStmt(stmt *ast.AssignStmt) {
	errCount := 0
	for _, expr := range stmt.Lhs {
		exprType := fc.pass.TypesInfo.TypeOf(expr)
		if exprType != nil && isError(exprType) {
			errCount++
			if errCount > 1 {
				fc.pass.Reportf(stmt.Pos(), "multiple errors")
				return
			}
		}
	}

	fc.checkAssignment(stmt.Lhs, stmt.Rhs)
}

func (fc *fileChecker) checkAssignment(lsh, rsh []ast.Expr) {
	if len(lsh) == len(rsh) {
		for i := range lsh {
			if fc.shouldWrap(rsh[i]) {
				fc.assignedErrorDst = ast.Unparen(lsh[i])
				fc.assignedErrorSrc = ast.Unparen(rsh[i])
				return
			}
		}
	} else {
		call, ok := ast.Unparen(rsh[0]).(*ast.CallExpr)
		if !ok {
			return
		}

		if fc.shouldWrap(call) {
			fc.assignedErrorDst = ast.Unparen(lsh[fc.errorReturnIndex(call)])
			fc.assignedErrorSrc = call
		}
	}
}

func (fc *fileChecker) checkAssignmentWrapping(node ast.Node) {
	if fc.assignedErrorDst != nil && !fc.assignmentTracker.isInNode() {
		assignedErrorDst := fc.assignedErrorDst
		assignedErrorSrc := fc.assignedErrorSrc
		fc.assignedErrorDst = nil
		fc.assignedErrorSrc = nil

		assignStmt, ok := node.(*ast.AssignStmt)
		if !ok {
			fc.report(assignedErrorSrc)
			return
		}

		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			fc.report(assignedErrorSrc)
			return
		}

		call, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			fc.report(assignedErrorSrc)
			return
		}

		if !fc.isWrapCall(call) {
			fc.report(assignedErrorSrc)
			return
		}

		if !areExprsEqual(assignStmt.Lhs[0], assignedErrorDst) {
			fc.report(assignedErrorSrc)
			return
		}

		if len(call.Args) != 1 {
			fc.report(assignedErrorSrc)
			return
		}

		if !areExprsEqual(call.Args[0], assignedErrorDst) {
			fc.report(assignedErrorSrc)
			return
		}
	}
}

func (fc *fileChecker) isWrapCall(call *ast.CallExpr) bool {
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

	obj := fc.pass.TypesInfo.Uses[ident]

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkg.Imported().Path() == "github.com/tbeati/stacked"
}

func (fc *fileChecker) isInternalCall(call *ast.CallExpr) bool {
	if fc.isTypeConversion(call) {
		return false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	pkg := fc.pass.TypesInfo.Uses[selector.Sel].Pkg()
	if pkg == nil {
		return false
	}

	if fc.config.isGeneratedPackage(pkg.Path()) {
		return false
	}

	return strings.HasPrefix(pkg.Path(), fc.pass.Module.Path)
}

func (fc *fileChecker) isTypeConversion(call *ast.CallExpr) bool {
	var obj types.Object

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		obj = fc.pass.TypesInfo.Uses[fun]
	case *ast.SelectorExpr:
		obj = fc.pass.TypesInfo.Uses[fun.Sel]
	}

	if obj == nil {
		return false
	}

	_, isTypeName := obj.(*types.TypeName)
	return isTypeName
}

func (fc *fileChecker) errorReturnIndex(call *ast.CallExpr) int {
	switch returnType := fc.pass.TypesInfo.TypeOf(call).(type) {
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

func (fc *fileChecker) returnsError(call *ast.CallExpr) bool {
	return fc.errorReturnIndex(call) >= 0
}
