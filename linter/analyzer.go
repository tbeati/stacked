package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const stackedImportPath = "github.com/tbeati/stacked"

type Config struct {
	PackagesTreatedAsExternal []string `json:"packages-treated-as-external"`
	IgnoredFunctions          []string `json:"ignored-functions"`
}

func (c *Config) isPackageTreatedAsExternal(pkg string) bool {
	for _, genPkg := range c.PackagesTreatedAsExternal {
		if strings.HasPrefix(pkg, genPkg) {
			return true
		}
	}

	return false
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	if config == nil {
		config = &Config{}
	}

	config.IgnoredFunctions = append(config.IgnoredFunctions,
		"errors.Join",
		"errors.Unwrap",
	)

	return &analysis.Analyzer{
		Name: "stacked",
		Doc:  "check for error not wrapped with stacked",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			if config.isPackageTreatedAsExternal(pass.Pkg.Path()) {
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
	config       *Config
	pass         *analysis.Pass
	file         *ast.File
	ignoredLines map[int]struct{}

	functionTracker nodeTracker
}

func newFileChecker(config *Config, pass *analysis.Pass, file *ast.File) *fileChecker {
	ignoredLines := make(map[int]struct{})
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			directive, ok := ast.ParseDirective(comment.Pos(), comment.Text)
			if ok && directive.Tool == "stacked" && directive.Name == "disable" {
				line := pass.Fset.Position(comment.Pos()).Line
				ignoredLines[line] = struct{}{}
			}
		}
	}

	return &fileChecker{
		config:          config,
		pass:            pass,
		file:            file,
		ignoredLines:    ignoredLines,
		functionTracker: newNodeTracker(),
	}
}

func (fc *fileChecker) check() {
	ast.Inspect(fc.file, func(node ast.Node) bool {
		fc.functionTracker.depthFirstSearchStep(node)
		fc.trackTopLevelFunctionDeclaration(node)

		if !fc.isInFunction() {
			return true
		}

		switch node := node.(type) {
		case *ast.DeclStmt:
		case *ast.GenDecl:
			fc.checkGenDecl(node)
		case *ast.AssignStmt:
			fc.checkAssignStmt(node)
		case *ast.ReturnStmt:
			for _, result := range node.Results {
				if fc.shouldWrap(result) {
					fc.report(result)
				}
			}
		case *ast.CallExpr:
			if fc.isWrapCall(node) {
				return true
			}

			isErrorIsOrAs := fc.isErrorIsOrAs(node)

			for i, arg := range node.Args {
				if isErrorIsOrAs && i == 1 {
					continue
				}

				if fc.shouldWrap(arg) {
					fc.report(arg)
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
		case *ast.SendStmt:
			if fc.shouldWrap(node.Value) {
				fc.report(node.Value)
			}
		case *ast.RangeStmt:
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

func (fc *fileChecker) checkGenDecl(stmt *ast.GenDecl) {
	if stmt.Tok != token.VAR {
		return
	}

SpecLoop:
	for _, spec := range stmt.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if len(valueSpec.Values) == 0 {
			continue
		}

		errCount := 0
		for _, name := range valueSpec.Names {
			if isError(fc.pass.TypesInfo.TypeOf(name)) {
				errCount++
				if errCount > 1 {
					fc.pass.Reportf(stmt.Pos(), "multiple errors") //TODO: improve message
					continue SpecLoop
				}
			}
		}

		lsh := make([]ast.Expr, 0, len(valueSpec.Names))
		for _, ident := range valueSpec.Names {
			lsh = append(lsh, ident)
		}

		fc.checkAssignment(lsh, valueSpec.Values)
	}
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
			if !isBlankIdent(lsh[i]) && fc.shouldWrap(rsh[i]) {
				fc.report(ast.Unparen(rsh[i]))
				return
			}
		}
	} else {
		call, ok := ast.Unparen(rsh[0]).(*ast.CallExpr)
		if !ok {
			return
		}

		if fc.shouldWrap(call) {
			assignedErrorDst := ast.Unparen(lsh[fc.errorReturnIndex(call)])
			if !isBlankIdent(assignedErrorDst) {
				fc.report(call)
			}
		}
	}
}

func (fc *fileChecker) report(expr ast.Expr) {
	var msg string
	var valueCount = 1

	expr = ast.Unparen(expr)
	switch expr := expr.(type) {
	case *ast.Ident:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	case *ast.SelectorExpr:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	case *ast.CompositeLit:
		msg = fmt.Sprintf("%s literal is not wrapped with stacked", exprToString(expr.Type))
	case *ast.UnaryExpr:
		switch subExpr := expr.X.(type) {
		case *ast.CompositeLit:
			msg = fmt.Sprintf("%s literal is not wrapped with stacked", exprToString(subExpr.Type))
		default:
			if expr.Op == token.ARROW {
				msg = fmt.Sprintf("error received from %s is not wrapped with stacked", exprToString(expr.X))
			} else {
				msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
			}
		}
	case *ast.StarExpr:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	case *ast.CallExpr:
		if fc.isTypeConversion(expr) {
			msg = fmt.Sprintf("value converted to error type %s is not wrapped with stacked", exprToString(expr.Fun))
		} else {
			msg = fmt.Sprintf("error returned by %s is not wrapped with stacked", exprToString(expr.Fun))
			valueCount = fc.returnValueCount(expr)
		}
	}

	var suggestedFixes []analysis.SuggestedFix
	if valueCount <= 5 {
		wrapValueCountString := ""
		if valueCount > 1 {
			wrapValueCountString = strconv.Itoa(valueCount)
		}

		suggestedFixes = []analysis.SuggestedFix{{
			Message: "", // TODO: add message
			TextEdits: []analysis.TextEdit{{
				Pos:     expr.Pos(),
				End:     expr.End(),
				NewText: []byte(fmt.Sprintf("stacked.Wrap%s(%s)", wrapValueCountString, exprToString(expr))),
			}},
		}}

		addMissingImport := addMissingStackedImport(fc.file)
		if addMissingImport != nil {
			suggestedFixes[0].TextEdits = append(suggestedFixes[0].TextEdits, *addMissingImport)
		}
	}

	fc.pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func addMissingStackedImport(file *ast.File) *analysis.TextEdit {
	for _, imp := range file.Imports {
		importPath, err := strconv.Unquote(imp.Path.Value)
		if err == nil && importPath == stackedImportPath {
			return nil
		}
	}

	importPos := file.Name.End()
	newText := fmt.Sprintf("\nimport \"%s\"\n", stackedImportPath)

	return &analysis.TextEdit{
		Pos:     importPos,
		End:     importPos,
		NewText: []byte(newText),
	}
}

func (fc *fileChecker) shouldWrap(expr ast.Expr) bool {
	if fc.isIgnoredLine(expr) {
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
	case *ast.UnaryExpr:
		return fc.shouldWrapUnary(expr)
	case *ast.StarExpr:
	case *ast.CallExpr:
		return fc.shouldWrapCall(expr)
	}

	return false
}

func (fc *fileChecker) isIgnoredLine(expr ast.Expr) bool {
	line := fc.pass.Fset.Position(expr.Pos()).Line
	_, isIgnored := fc.ignoredLines[line]
	return isIgnored
}

func (fc *fileChecker) shouldWrapIdent(ident *ast.Ident) bool {
	obj := fc.pass.TypesInfo.ObjectOf(ident)

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return isError(variable.Type()) && variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()
}

func (fc *fileChecker) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	obj := fc.pass.TypesInfo.ObjectOf(expr.Sel)

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return isError(variable.Type()) && variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()
}

func (fc *fileChecker) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return isError(fc.pass.TypesInfo.TypeOf(lit))
}

func (fc *fileChecker) shouldWrapUnary(expr *ast.UnaryExpr) bool {
	return isError(fc.pass.TypesInfo.TypeOf(expr))
}

func (fc *fileChecker) shouldWrapCall(call *ast.CallExpr) bool {
	if fc.isWrapCall(call) {
		return false
	}

	if fc.isTypeConversion(call) {
		return fc.returnsError(call)
	}

	if fc.isFunctionLiteral(call) {
		return fc.returnsError(call)
	}

	if fc.isInternalCall(call) {
		return false
	}

	if fc.isIgnoredCall(call) {
		return false
	}

	return fc.returnsError(call)
}

func (fc *fileChecker) isWrapCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	wrapFunctions := map[string]struct{}{
		"Wrap":      {},
		"Wrap2":     {},
		"Wrap3":     {},
		"Wrap4":     {},
		"Wrap5":     {},
		"WrapSeq":   {},
		"WrapSeq2":  {},
		"WrapPull":  {},
		"WrapPull2": {},
	}
	_, isWrapFunction := wrapFunctions[selector.Sel.Name]
	if !isWrapFunction {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := fc.pass.TypesInfo.ObjectOf(ident)

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkg.Imported().Path() == stackedImportPath
}

func (fc *fileChecker) isTypeConversion(call *ast.CallExpr) bool {
	callType, ok := fc.pass.TypesInfo.Types[call.Fun]
	return ok && callType.IsType()
}

func (fc *fileChecker) isFunctionLiteral(call *ast.CallExpr) bool {
	_, isFuncLit := call.Fun.(*ast.FuncLit)
	return isFuncLit
}

func (fc *fileChecker) isInternalCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		obj := fc.pass.TypesInfo.ObjectOf(fun)

		_, isFunc := obj.(*types.Func)
		if !isFunc {
			return false
		}

		return obj.Pkg() != nil
	case *ast.SelectorExpr:
		obj := fc.pass.TypesInfo.ObjectOf(fun.Sel)

		sel := fc.pass.TypesInfo.Selections[fun]
		if sel != nil {
			if sel.Kind() == types.FieldVal {
				return false
			}

			recvType := sel.Recv()
			_, isInterface := recvType.Underlying().(*types.Interface)
			if isInterface {
				return false
			}
		} else {
			_, isFunc := obj.(*types.Func)
			if !isFunc {
				return false
			}
		}

		pkg := obj.Pkg()
		if pkg == nil {
			return false
		}

		if fc.config.isPackageTreatedAsExternal(pkg.Path()) {
			return false
		}

		return strings.HasPrefix(pkg.Path(), fc.pass.Module.Path)
	}

	return false
}

func (fc *fileChecker) isIgnoredCall(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	sel := fc.pass.TypesInfo.Selections[selector]
	if sel != nil {
		if sel.Kind() == types.FieldVal {
			return false
		}

		methodPath := strings.TrimPrefix(sel.Recv().String(), "*") + "." + sel.Obj().Name()
		for _, ignoredFunction := range fc.config.IgnoredFunctions {
			if methodPath == ignoredFunction {
				return true
			}
		}
	} else {
		obj := fc.pass.TypesInfo.ObjectOf(selector.Sel)

		_, isFunc := obj.(*types.Func)
		if !isFunc {
			return false
		}

		pkg := obj.Pkg()
		if pkg == nil {
			return false
		}

		functionPath := pkg.Path() + "." + selector.Sel.Name
		for _, ignoredFunction := range fc.config.IgnoredFunctions {
			if functionPath == ignoredFunction {
				return true
			}
		}

		if pkg.Path() == "fmt" && selector.Sel.Name == "Errorf" {
			formatString, ok := call.Args[0].(*ast.BasicLit)
			if !ok || formatString.Kind != token.STRING {
				return false
			}

			if strings.Contains(formatString.Value, "%w") {
				return true
			}
		}
	}

	return false
}

func (fc *fileChecker) returnsError(call *ast.CallExpr) bool {
	return fc.errorReturnIndex(call) >= 0
}

func (fc *fileChecker) errorReturnIndex(call *ast.CallExpr) int {
	returnType := fc.pass.TypesInfo.TypeOf(call)

	tuple, ok := returnType.(*types.Tuple)
	if ok {
		for i := range tuple.Len() {
			if isError(tuple.At(i).Type()) {
				return i
			}
		}
	}

	if isError(returnType) {
		return 0
	}

	return -1
}

func (fc *fileChecker) returnValueCount(call *ast.CallExpr) int {
	callType := fc.pass.TypesInfo.TypeOf(call)

	tuple, ok := callType.(*types.Tuple)
	if ok {
		return tuple.Len()
	}

	return 1
}

func (fc *fileChecker) isErrorIsOrAs(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if selector.Sel.Name != "Is" && selector.Sel.Name != "As" {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := fc.pass.TypesInfo.ObjectOf(ident)

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkg.Imported().Path() == "errors"
}

func (fc *fileChecker) isErrorIterator(expr ast.Expr) (bool, int) {
	exprType := fc.pass.TypesInfo.TypeOf(expr)

	sig, ok := types.Unalias(exprType).Underlying().(*types.Signature)
	if !ok {
		return false, 0
	}

	params := sig.Params()
	results := sig.Results()

	if params.Len() != 1 || results.Len() != 0 {
		return false, 0
	}

	yieldSig, ok := types.Unalias(params.At(0).Type()).Underlying().(*types.Signature)
	if !ok {
		return false, 0
	}
	_ = yieldSig

	return false, 0
}
