package linter

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/printer"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const stackedImportPath = "github.com/tbeati/stacked"

type Config struct {
	PackagesTreatedAsExternal []string           `json:"packages-treated-as-external"`
	IgnoredFunctions          []string           `json:"ignored-functions"`
	CheckFunctionArguments    []FunctionArgument `json:"check-function-arguments"`

	ignoredFunctionsMap       map[string]struct{}
	checkFunctionArgumentsMap map[string]int
}

type FunctionArgument struct {
	Function string
	Argument int
}

func (c *Config) init() {
	c.ignoredFunctionsMap = make(map[string]struct{}, len(c.IgnoredFunctions))
	for _, fun := range c.IgnoredFunctions {
		c.ignoredFunctionsMap[fun] = struct{}{}
	}

	c.checkFunctionArgumentsMap = make(map[string]int, len(c.CheckFunctionArguments))
	for _, checkFun := range c.CheckFunctionArguments {
		c.checkFunctionArgumentsMap[checkFun.Function] = checkFun.Argument
	}
}

func (c *Config) isPackageTreatedAsExternal(pkg string) bool {
	for _, genPkg := range c.PackagesTreatedAsExternal {
		if strings.HasPrefix(pkg, genPkg) {
			return true
		}
	}

	return false
}

func (c *Config) isIgnoredFunction(function string) bool {
	_, found := c.ignoredFunctionsMap[function]
	return found
}

func (c *Config) isCheckFunction(function string) bool {
	_, exists := c.checkFunctionArgumentsMap[function]
	return exists
}

func (c *Config) checkFunctionArgument(function string) int {
	return c.checkFunctionArgumentsMap[function]
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	if config == nil {
		config = &Config{}
	}

	config.IgnoredFunctions = append(config.IgnoredFunctions,
		"errors.Join",
		"errors.Unwrap",
	)

	config.CheckFunctionArguments = append(config.CheckFunctionArguments,
		FunctionArgument{
			Function: "errors.Is",
			Argument: 2,
		},
		FunctionArgument{
			Function: "errors.As",
			Argument: 2,
		},
	)

	config.init()

	return &analysis.Analyzer{
		Name:             "stacked",
		Doc:              "check for error not wrapped with stacked",
		URL:              "https://github.com/tbeati/stacked",
		RunDespiteErrors: false,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			if config.isPackageTreatedAsExternal(pass.Pkg.Path()) {
				return nil, nil
			}

			newAnalyzer(config, pass).analyze()

			return nil, nil
		},
	}
}

type analyzer struct {
	config       *Config
	pass         *analysis.Pass
	ignoredLines map[*ast.File]map[int]struct{}

	stack []ast.Node
}

func newAnalyzer(config *Config, pass *analysis.Pass) *analyzer {
	ignoredLines := make(map[*ast.File]map[int]struct{})
	for _, file := range pass.Files {
		ignoredLines[file] = make(map[int]struct{})
		for _, commentGroup := range file.Comments {
			for _, comment := range commentGroup.List {
				if strings.Contains(comment.Text, "//stacked:disable") {
					line := pass.Fset.Position(comment.Pos()).Line
					ignoredLines[file][line] = struct{}{}
				}
			}
		}
	}

	return &analyzer{
		config:       config,
		pass:         pass,
		ignoredLines: ignoredLines,
	}
}

func (a *analyzer) analyze() {
	astInspector := a.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	astInspector.WithStack(nil, func(node ast.Node, push bool, stack []ast.Node) bool {
		a.stack = stack

		if !push {
			return true
		}

		if a.enclosingFunctionSignature() == nil {
			return true
		}

		switch node := node.(type) {
		case *ast.GenDecl:
			a.checkGenDecl(node)
		case *ast.AssignStmt:
			a.checkAssignStmt(node)
		case *ast.ReturnStmt:
			for i, result := range node.Results {
				if a.shouldWrap(result) {
					a.report(result, a.isErrorExpectedInReturn(i, len(node.Results)))
				}
			}
		case *ast.CallExpr:
			if a.isStackedWrapCall(node) {
				return true
			}

			isTypeConversion := a.isTypeConversion(node)
			isErrorCheckCall, argumentIndex := a.isErrorCheckCall(node)

			for i, arg := range node.Args {
				if isErrorCheckCall && i == argumentIndex {
					continue
				}

				if isTypeConversion && a.isConst(arg) {
					continue
				}

				if a.shouldWrap(arg) {
					a.report(arg, a.isErrorExpectedInCallArgs(node, i, len(node.Args)))
				}
			}
		case *ast.CompositeLit:
			for i, elt := range node.Elts {
				switch unparenElt := ast.Unparen(elt).(type) {
				case *ast.KeyValueExpr:
					if a.shouldWrap(unparenElt.Value) {
						a.report(unparenElt.Value, a.isErrorExpectedInLit(node, unparenElt, i))
					}
				default:
					if a.shouldWrap(elt) {
						a.report(elt, a.isErrorExpectedInLit(node, unparenElt, i))
					}
				}
			}
		case *ast.SendStmt:
			if a.shouldWrap(node.Value) {
				a.report(node.Value, a.isErrorExpectedInSend(node))
			}
		case *ast.RangeStmt:
			shouldWrap, autoFixable, returnValueCount := a.shouldWrapIterator(node.X)
			if shouldWrap {
				a.reportIterator(node.X, autoFixable, returnValueCount)
			}
		}

		return true
	})
}

func (a *analyzer) currentFile() *ast.File {
	return a.stack[0].(*ast.File)
}

func (a *analyzer) enclosingFunctionSignature() *types.Signature {
	for i := len(a.stack) - 1; i >= 0; i-- {
		switch fun := a.stack[i].(type) {
		case *ast.FuncDecl:
			return a.pass.TypesInfo.ObjectOf(fun.Name).Type().Underlying().(*types.Signature)
		case *ast.FuncLit:
			return a.pass.TypesInfo.TypeOf(fun.Type).Underlying().(*types.Signature)
		}
	}

	return nil
}

func (a *analyzer) isErrorCheckCall(call *ast.CallExpr) (bool, int) {
	isErrorCheckCall, functionName := a.isCalledFunctionInSet(call, a.config.isCheckFunction)
	if !isErrorCheckCall {
		return false, 0
	}

	argumentIndex := a.config.checkFunctionArgument(functionName) - 1

	return true, argumentIndex
}

func (a *analyzer) checkGenDecl(stmt *ast.GenDecl) {
	if stmt.Tok != token.VAR {
		return
	}

	for _, spec := range stmt.Specs {
		valueSpec := spec.(*ast.ValueSpec)

		if len(valueSpec.Values) == 0 {
			continue
		}

		lhs := make([]ast.Expr, 0, len(valueSpec.Names))
		for _, ident := range valueSpec.Names {
			lhs = append(lhs, ident)
		}

		if a.multipleErrorsInAssigment(lhs, stmt.Pos()) {
			continue
		}

		a.checkAssignment(lhs, valueSpec.Values)
	}
}

func (a *analyzer) checkAssignStmt(stmt *ast.AssignStmt) {
	if a.multipleErrorsInAssigment(stmt.Lhs, stmt.Pos()) {
		return
	}

	a.checkAssignment(stmt.Lhs, stmt.Rhs)
}

func (a *analyzer) multipleErrorsInAssigment(lhs []ast.Expr, assignPos token.Pos) bool {
	errCount := 0
	for _, expr := range lhs {
		exprType := a.pass.TypesInfo.TypeOf(expr)
		if exprType != nil && implementsError(exprType) {
			errCount++
			if errCount > 1 {
				a.pass.Reportf(assignPos, "assignment to multiple error variables")
				return true
			}
		}
	}

	return false
}

func (a *analyzer) checkAssignment(lhs, rhs []ast.Expr) {
	if len(lhs) == len(rhs) {
		for i := range lhs {
			if !isBlankIdent(lhs[i]) && a.shouldWrap(rhs[i]) {
				a.report(rhs[i], isError(a.pass.TypesInfo.TypeOf(lhs[i])))
				return
			}
		}
	} else {
		switch unparenRhs := ast.Unparen(rhs[0]).(type) {
		case *ast.CallExpr:
			if a.shouldWrap(rhs[0]) {
				assignedErrorVariable := lhs[a.errorReturnIndex(unparenRhs)]
				if !isBlankIdent(assignedErrorVariable) {
					a.report(rhs[0], isError(a.pass.TypesInfo.TypeOf(assignedErrorVariable)))
				}
			}
		case *ast.UnaryExpr:
			if a.shouldWrap(rhs[0]) && !isBlankIdent(lhs[0]) {
				a.report(rhs[0], false)
			}
		}
	}
}

func (a *analyzer) report(expr ast.Expr, autoFixable bool) {
	var msg string
	valueCount := 1
	isIteratorPull := false

	switch expr := ast.Unparen(expr).(type) {
	case *ast.BasicLit:
		msg = fmt.Sprintf("basic literal %s is not wrapped with stacked", a.exprToString(expr))
	case *ast.Ident:
		obj := a.pass.TypesInfo.ObjectOf(expr)
		if obj.Parent() == types.Universe && (obj.Name() == "true" || obj.Name() == "false") {
			msg = fmt.Sprintf("basic literal %s is not wrapped with stacked", a.exprToString(expr))
		} else {
			msg = fmt.Sprintf("%s is not wrapped with stacked", a.exprToString(expr))
		}
	case *ast.SelectorExpr:
		msg = fmt.Sprintf("%s is not wrapped with stacked", a.exprToString(expr))
	case *ast.CompositeLit:
		exprType := a.pass.TypesInfo.TypeOf(expr)
		msg = fmt.Sprintf("%s literal is not wrapped with stacked", typeToString(exprType, a.pass.Pkg))
	case *ast.UnaryExpr:
		switch subExpr := ast.Unparen(expr.X).(type) {
		case *ast.CompositeLit:
			exprType := a.pass.TypesInfo.TypeOf(subExpr)
			msg = fmt.Sprintf("%s literal is not wrapped with stacked", typeToString(exprType, a.pass.Pkg))
		default:
			if expr.Op == token.ARROW {
				msg = fmt.Sprintf("error received from %s is not wrapped with stacked", a.exprToString(subExpr))
			}
		}
	case *ast.CallExpr:
		fun := ast.Unparen(expr.Fun)
		if a.isTypeConversion(expr) {
			msg = fmt.Sprintf("value converted to error type %s is not wrapped with stacked", a.exprToString(fun))
		} else {
			msg = fmt.Sprintf("error returned by %s is not wrapped with stacked", a.exprToString(fun))
			var isCallExprAutoFixable bool
			isCallExprAutoFixable, isIteratorPull, valueCount = a.isCallExprAutoFixable(expr)
			autoFixable = autoFixable && isCallExprAutoFixable
		}
	}

	var suggestedFixes []analysis.SuggestedFix
	if autoFixable {
		wrapPull := ""
		if isIteratorPull {
			wrapPull = "Pull"
		}

		wrapValueCountString := ""
		if valueCount > 1 {
			wrapValueCountString = strconv.Itoa(valueCount)
		}

		wrapFuncName := fmt.Sprintf("stacked.Wrap%s%s", wrapPull, wrapValueCountString)

		suggestedFixes = a.suggestedFixes(expr, wrapFuncName)
	}

	a.pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func (a *analyzer) reportIterator(expr ast.Expr, autoFixable bool, returnValueCount int) {
	var msg string

	switch expr := ast.Unparen(expr).(type) {
	case *ast.FuncLit:
		msg = fmt.Sprintf("iterator literal is not wrapped with stacked")
	case *ast.UnaryExpr:
		if expr.Op == token.ARROW {
			msg = fmt.Sprintf("iterator received from %s is not wrapped with stacked", a.exprToString(ast.Unparen(expr.X)))
		} else {
			msg = fmt.Sprintf("%s is not wrapped with stacked", a.exprToString(expr))
		}
	case *ast.CallExpr:
		fun := ast.Unparen(expr.Fun)
		if a.isTypeConversion(expr) {
			msg = fmt.Sprintf("value converted to iterator type %s is not wrapped with stacked", a.exprToString(fun))
		} else {
			msg = fmt.Sprintf("iterator returned by %s is not wrapped with stacked", a.exprToString(fun))
		}
	default:
		msg = fmt.Sprintf("%s is not wrapped with stacked", a.exprToString(expr))
	}

	var suggestedFixes []analysis.SuggestedFix
	if autoFixable && returnValueCount <= 2 {
		wrapValueCountString := ""
		if returnValueCount > 1 {
			wrapValueCountString = strconv.Itoa(returnValueCount)
		}

		wrapFuncName := fmt.Sprintf("stacked.WrapSeq%s", wrapValueCountString)

		suggestedFixes = a.suggestedFixes(expr, wrapFuncName)
	}

	a.pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func (a *analyzer) suggestedFixes(expr ast.Expr, wrapFuncName string) []analysis.SuggestedFix {
	suggestedFixes := []analysis.SuggestedFix{{
		Message: fmt.Sprintf("Wrap with %s", wrapFuncName),
		TextEdits: []analysis.TextEdit{{
			Pos:     expr.Pos(),
			End:     expr.End(),
			NewText: []byte(fmt.Sprintf("%s(%s)", wrapFuncName, a.exprToString(expr))),
		}},
	}}

	addMissingImport := addMissingStackedImport(a.currentFile())
	if addMissingImport != nil {
		suggestedFixes[0].TextEdits = append(suggestedFixes[0].TextEdits, *addMissingImport)
	}

	return suggestedFixes
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

func (a *analyzer) exprToString(expr ast.Expr) string {
	var s strings.Builder
	err := printer.Fprint(&s, a.pass.Fset, expr)
	if err != nil {
		panic(err)
	}
	return s.String()
}

func (a *analyzer) shouldWrap(expr ast.Expr) bool {
	if a.isIgnoredLine(expr) {
		return false
	}

	expr = ast.Unparen(expr)
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return a.shouldWrapBasicLit(expr)
	case *ast.Ident:
		return a.shouldWrapIdent(expr)
	case *ast.SelectorExpr:
		return a.shouldWrapSelector(expr)
	case *ast.CompositeLit:
		return a.shouldWrapCompositeLit(expr)
	case *ast.UnaryExpr:
		return a.shouldWrapUnary(expr)
	case *ast.CallExpr:
		return a.shouldWrapCall(expr)
	}

	return false
}

func (a *analyzer) isIgnoredLine(expr ast.Expr) bool {
	line := a.pass.Fset.Position(expr.Pos()).Line
	_, isIgnored := a.ignoredLines[a.currentFile()][line]
	return isIgnored
}

func (a *analyzer) shouldWrapBasicLit(expr *ast.BasicLit) bool {
	return implementsError(a.pass.TypesInfo.TypeOf(expr))
}

func (a *analyzer) shouldWrapIdent(ident *ast.Ident) bool {
	if !implementsError(a.pass.TypesInfo.TypeOf(ident)) {
		return false
	}

	obj := a.pass.TypesInfo.ObjectOf(ident)

	switch obj.(type) {
	case *types.Const:
		return true
	case *types.Var:
		return obj.Pkg() != nil && obj.Parent() == obj.Pkg().Scope()
	}

	return false
}

func (a *analyzer) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	if !implementsError(a.pass.TypesInfo.TypeOf(expr)) {
		return false
	}

	obj := a.pass.TypesInfo.ObjectOf(expr.Sel)

	switch obj.(type) {
	case *types.Const:
		return true
	case *types.Var:
		return obj.Pkg() != nil && obj.Parent() == obj.Pkg().Scope()
	}

	return false
}

func (a *analyzer) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return implementsError(a.pass.TypesInfo.TypeOf(lit))
}

func (a *analyzer) shouldWrapUnary(expr *ast.UnaryExpr) bool {
	exprType := a.pass.TypesInfo.TypeOf(expr)

	switch ast.Unparen(expr.X).(type) {
	case *ast.CompositeLit:
		return implementsError(exprType)
	default:
		if expr.Op == token.ARROW {
			tuple, isTuple := exprType.Underlying().(*types.Tuple)
			if isTuple {
				return implementsError(tuple.At(0).Type())
			}

			return implementsError(exprType)
		}
	}

	return false
}

func (a *analyzer) shouldWrapCall(call *ast.CallExpr) bool {
	if a.isStackedWrapCall(call) {
		return false
	}

	if a.isTypeConversion(call) {
		return a.returnsError(call) && (a.isConst(call.Args[0]) || !implementsError(a.pass.TypesInfo.TypeOf(call.Args[0])))
	}

	if a.isInternalCall(call) {
		return false
	}

	if a.isIgnoredCall(call) {
		return false
	}

	return a.returnsError(call)
}

func (a *analyzer) isStackedWrapCall(call *ast.CallExpr) bool {
	var ident *ast.Ident
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.SelectorExpr:
		ident = fun.Sel
	case *ast.Ident:
		ident = fun
	}

	if ident == nil {
		return false
	}

	switch ident.Name {
	case "Wrap", "Wrap2", "Wrap3", "Wrap4", "Wrap5", "WrapSeq", "WrapSeq2", "WrapPull", "WrapPull2":
	default:
		return false
	}

	obj := a.pass.TypesInfo.ObjectOf(ident)

	_, isFunc := obj.(*types.Func)
	if !isFunc {
		return false
	}

	return obj.Pkg() != nil && obj.Pkg().Path() == stackedImportPath
}

func (a *analyzer) isInternalCall(call *ast.CallExpr) bool {
	var ident *ast.Ident
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.FuncLit:
		return true
	case *ast.SelectorExpr:
		sel := a.pass.TypesInfo.Selections[fun]
		if sel != nil {
			if sel.Kind() == types.FieldVal {
				return false
			}

			_, isInterface := sel.Recv().Underlying().(*types.Interface)
			if isInterface {
				return false
			}
		}

		ident = fun.Sel
	case *ast.Ident:
		ident = fun
	}

	if ident == nil {
		return false
	}

	obj := a.pass.TypesInfo.ObjectOf(ident)

	_, isFunc := obj.(*types.Func)
	if !isFunc {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	if a.config.isPackageTreatedAsExternal(pkg.Path()) {
		return false
	}

	return strings.HasPrefix(pkg.Path(), a.pass.Module.Path)
}

func (a *analyzer) isIgnoredCall(call *ast.CallExpr) bool {
	isIgnoredCall, _ := a.isCalledFunctionInSet(call, a.config.isIgnoredFunction)
	return isIgnoredCall || a.isFmtErrorfWithW(call)
}

func (a *analyzer) isFmtErrorfWithW(call *ast.CallExpr) bool {
	var ident *ast.Ident
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		ident = fun
	case *ast.SelectorExpr:
		ident = fun.Sel
	}

	if ident == nil {
		return false
	}

	obj := a.pass.TypesInfo.ObjectOf(ident)

	if obj.Pkg() == nil || obj.Pkg().Path() != "fmt" || obj.Name() != "Errorf" {
		return false
	}

	if len(call.Args) == 0 {
		return false
	}

	formatArg := call.Args[0]
	tv := a.pass.TypesInfo.Types[formatArg]
	if tv.Value == nil || tv.Value.Kind() != constant.String {
		return false
	}

	formatString := constant.StringVal(tv.Value)

	return containsWVerb(formatString)
}

func (a *analyzer) returnsError(call *ast.CallExpr) bool {
	return a.errorReturnIndex(call) >= 0
}

func (a *analyzer) errorReturnIndex(call *ast.CallExpr) int {
	returnType := a.pass.TypesInfo.TypeOf(call)

	tuple, ok := returnType.Underlying().(*types.Tuple)
	if ok {
		for i := range tuple.Len() {
			if implementsError(tuple.At(i).Type()) {
				return i
			}
		}
	}

	if implementsError(returnType) {
		return 0
	}

	return -1
}

func (a *analyzer) isTypeConversion(call *ast.CallExpr) bool {
	callType, ok := a.pass.TypesInfo.Types[call.Fun]
	return ok && callType.IsType()
}

func (a *analyzer) isConst(expr ast.Expr) bool {
	return a.pass.TypesInfo.Types[expr].Value != nil
}

func (a *analyzer) isCallExprAutoFixable(call *ast.CallExpr) (bool, bool, int) {
	callType := a.pass.TypesInfo.TypeOf(call)

	tuple, ok := callType.Underlying().(*types.Tuple)
	if ok {
		if a.isIteratorPull(call) {
			return true, true, tuple.Len() - 1
		}

		errorReturnIndex := a.errorReturnIndex(call)
		returnValueCount := tuple.Len()
		return returnValueCount <= 5 && errorReturnIndex == returnValueCount-1, false, returnValueCount
	}

	return true, false, 1
}

func (a *analyzer) isIteratorPull(call *ast.CallExpr) bool {
	funType := a.pass.TypesInfo.TypeOf(call.Fun)

	sig, ok := funType.Underlying().(*types.Signature)
	if !ok {
		return false
	}

	if sig.Params().Len() > 0 {
		return false
	}

	results := sig.Results()
	if results == nil {
		return false
	}

	switch results.Len() {
	case 2:
		return implementsError(results.At(0).Type()) && isBool(results.At(1).Type())
	case 3:
		return implementsError(results.At(1).Type()) && isBool(results.At(2).Type())
	}

	return false
}

func (a *analyzer) isCalledFunctionInSet(call *ast.CallExpr, isFunctionInSet func(function string) bool) (bool, string) {
	var ident *ast.Ident
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		ident = fun
	case *ast.SelectorExpr:
		ident = fun.Sel
	}

	if ident == nil {
		return false, ""
	}

	obj := a.pass.TypesInfo.ObjectOf(ident)

	fun, ok := obj.(*types.Func)
	if !ok {
		return false, ""
	}

	pkg := fun.Pkg()
	if pkg == nil {
		return false, ""
	}

	fullName := pkg.Path()

	sig, ok := fun.Type().(*types.Signature)
	if ok && sig.Recv() != nil {
		recvType := sig.Recv().Type()

		ptr, isPtr := recvType.(*types.Pointer)
		if isPtr {
			recvType = ptr.Elem()
		}

		named, isNamed := recvType.(*types.Named)
		if !isNamed {
			return false, ""
		}
		fullName += "." + named.Obj().Name()
	}

	fullName += "." + fun.Name()

	return isFunctionInSet(fullName), fullName
}

func (a *analyzer) shouldWrapIterator(expr ast.Expr) (bool, bool, int) {
	callExpr, isCall := expr.(*ast.CallExpr)
	if isCall {
		if a.isStackedWrapCall(callExpr) {
			return false, false, 0
		}
	}

	sig, ok := a.pass.TypesInfo.TypeOf(expr).Underlying().(*types.Signature)
	if !ok {
		return false, false, 0
	}

	if !(sig.Params().Len() == 1 && sig.Results().Len() == 0) {
		return false, false, 0
	}

	yieldSig, ok := sig.Params().At(0).Type().Underlying().(*types.Signature)
	if !ok {
		return false, false, 0
	}

	if yieldSig.Results().Len() != 1 {
		return false, false, 0
	}
	resType, ok := yieldSig.Results().At(0).Type().Underlying().(*types.Basic)
	if !ok || resType.Kind() != types.Bool {
		return false, false, 0
	}

	yieldParams := yieldSig.Params()
	switch yieldParams.Len() {
	case 1:
		yieldParamType := yieldParams.At(0).Type()
		return implementsError(yieldParamType), isError(yieldParamType), 1
	case 2:
		yieldParam0Type := yieldParams.At(0).Type()
		yieldParam1Type := yieldParams.At(1).Type()
		param0ImplementsError := implementsError(yieldParam0Type)
		param1ImplementsError := implementsError(yieldParam1Type)

		if param0ImplementsError && param1ImplementsError {
			a.pass.Reportf(expr.Pos(), "iterator yields multiple errors")
			return false, false, 0
		}

		if param0ImplementsError {
			return true, false, 2
		}

		if param1ImplementsError {
			return true, isError(yieldParam1Type), 2
		}
	}

	return false, false, 0
}

func (a *analyzer) isErrorExpectedInReturn(resultIndex int, returnedItemCount int) bool {
	results := a.enclosingFunctionSignature().Results()

	if returnedItemCount == 1 && results.Len() > 1 {
		if isError(results.At(results.Len() - 1).Type()) {
			return true
		}
		if isError(results.At(results.Len() - 2).Type()) {
			return true
		}
		return false
	}

	return isError(results.At(resultIndex).Type())
}

func (a *analyzer) isErrorExpectedInCallArgs(call *ast.CallExpr, argIndex, argCount int) bool {
	funType := a.pass.TypesInfo.TypeOf(call.Fun)

	sig, ok := funType.Underlying().(*types.Signature)
	if ok {
		params := sig.Params()

		if argCount == 1 && params.Len() > 1 {
			if sig.Variadic() && isError(params.At(params.Len()-1).Type().Underlying().(*types.Slice).Elem()) {
				return true
			}

			if isError(params.At(params.Len() - 1).Type()) {
				return true
			}

			if isError(params.At(params.Len() - 2).Type()) {
				return true
			}
		} else {
			var argType types.Type
			if sig.Variadic() && argIndex >= params.Len()-1 {
				lastParam := params.At(params.Len() - 1).Type()
				sliceType := lastParam.Underlying().(*types.Slice)
				argType = sliceType.Elem()
			} else if argIndex < params.Len() {
				argType = params.At(argIndex).Type()
			}

			return isError(argType)
		}
	} else if a.pass.TypesInfo.Types[call.Fun].IsType() {
		return isError(funType)
	}

	return false
}

func (a *analyzer) isErrorExpectedInLit(lit *ast.CompositeLit, elt ast.Expr, eltIndex int) bool {
	litType := a.pass.TypesInfo.TypeOf(lit)

	ptr, ok := litType.Underlying().(*types.Pointer)
	if ok {
		litType = ptr.Elem()
	}

	var expectedType types.Type

	switch litType := litType.Underlying().(type) {
	case *types.Struct:
		kv, ok := elt.(*ast.KeyValueExpr)
		if ok {
			ident := kv.Key.(*ast.Ident)
			expectedType = a.pass.TypesInfo.ObjectOf(ident).Type()
		} else {
			expectedType = litType.Field(eltIndex).Type()
		}
	case *types.Slice:
		expectedType = litType.Elem()
	case *types.Array:
		expectedType = litType.Elem()
	case *types.Map:
		expectedType = litType.Elem()
	}

	return isError(expectedType)
}

func (a *analyzer) isErrorExpectedInSend(send *ast.SendStmt) bool {
	chanType := a.pass.TypesInfo.TypeOf(send.Chan).Underlying().(*types.Chan)
	return isError(chanType.Elem())
}
