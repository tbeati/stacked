package linter

import (
	"fmt"
	"go/ast"
	"go/constant"
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
	CheckFunctionArguments    []FunctionArgument `json:"ignored-arguments"`
}

type FunctionArgument struct {
	Function string
	Argument int
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
	for _, ignoredFunction := range c.IgnoredFunctions {
		if function == ignoredFunction {
			return true
		}
	}

	return false
}

func (c *Config) isCheckFunction(function string) bool {
	for _, checkFunction := range c.CheckFunctionArguments {
		if function == checkFunction.Function {
			return true
		}
	}

	return false
}

func (c *Config) checkFunctionArgument(function string) int {
	for _, checkFunctionArgument := range c.CheckFunctionArguments {
		if function == checkFunctionArgument.Function {
			return checkFunctionArgument.Argument
		}
	}

	return -1
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

	return &analysis.Analyzer{
		Name:     "stacked",
		Doc:      "check for error not wrapped with stacked",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			if config.isPackageTreatedAsExternal(pass.Pkg.Path()) {
				return nil, nil
			}

			newAnalyzer(config, pass).check()

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
				directive, ok := ast.ParseDirective(comment.Pos(), comment.Text)
				if ok && directive.Tool == "stacked" && directive.Name == "disable" {
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

func (a *analyzer) check() {
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
			if a.isWrapCall(node) {
				return true
			}

			isErrorCheckCall, argumentIndex := a.isErrorCheckCall(node)

			for i, arg := range node.Args {
				if isErrorCheckCall && i == argumentIndex {
					continue
				}

				if a.shouldWrap(arg) {
					a.report(arg, a.isErrorExpectedInCallArgs(node, i))
				}
			}
		case *ast.CompositeLit:
			for i, elt := range node.Elts {
				switch elt := ast.Unparen(elt).(type) {
				case *ast.KeyValueExpr:
					if a.shouldWrap(elt.Value) {
						a.report(elt.Value, a.isErrorExpectedInLit(node, elt, i))
					}
				default:
					if a.shouldWrap(elt) {
						a.report(elt, a.isErrorExpectedInLit(node, elt, i))
					}
				}
			}
		case *ast.SendStmt:
			if a.shouldWrap(node.Value) {
				a.report(node.Value, a.isErrorExpectedInSend(node))
			}
		case *ast.RangeStmt:
			shouldWrap, valueCount := a.shouldWrapIterator(node.X)
			if shouldWrap {
				a.reportIterator(node.X, valueCount)
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
			return a.pass.TypesInfo.ObjectOf(fun.Name).Type().(*types.Signature)
		case *ast.FuncLit:
			return a.pass.TypesInfo.TypeOf(fun.Type).(*types.Signature)
		}
	}

	return nil
}

func (a *analyzer) checkGenDecl(stmt *ast.GenDecl) {
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
			if implementsError(a.pass.TypesInfo.TypeOf(name)) {
				errCount++
				if errCount > 1 {
					a.pass.Reportf(stmt.Pos(), "assignment to multiple error variables")
					continue SpecLoop
				}
			}
		}

		lsh := make([]ast.Expr, 0, len(valueSpec.Names))
		for _, ident := range valueSpec.Names {
			lsh = append(lsh, ident)
		}

		a.checkAssignment(lsh, valueSpec.Values)
	}
}

func (a *analyzer) checkAssignStmt(stmt *ast.AssignStmt) {
	errCount := 0
	for _, expr := range stmt.Lhs {
		exprType := a.pass.TypesInfo.TypeOf(expr)
		if exprType != nil && implementsError(exprType) {
			errCount++
			if errCount > 1 {
				a.pass.Reportf(stmt.Pos(), "assignment to multiple error variables")
				return
			}
		}
	}

	a.checkAssignment(stmt.Lhs, stmt.Rhs)
}

func (a *analyzer) checkAssignment(lsh, rsh []ast.Expr) {
	if len(lsh) == len(rsh) {
		for i := range lsh {
			if !isBlankIdent(lsh[i]) && a.shouldWrap(rsh[i]) {
				a.report(ast.Unparen(rsh[i]), isError(a.pass.TypesInfo.TypeOf(lsh[i])))
				return
			}
		}
	} else {
		call, ok := ast.Unparen(rsh[0]).(*ast.CallExpr)
		if !ok {
			return
		}

		if a.shouldWrap(call) {
			assignedErrorDst := ast.Unparen(lsh[a.errorReturnIndex(call)])
			if !isBlankIdent(assignedErrorDst) {
				a.report(call, isError(a.pass.TypesInfo.TypeOf(assignedErrorDst)))
			}
		}
	}
}

func (a *analyzer) report(expr ast.Expr, autoWrap bool) {
	var msg string
	valueCount := 1
	errorReturnIndex := 0
	isIteratorPull := false

	expr = ast.Unparen(expr)
	switch expr := expr.(type) {
	case *ast.Ident:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	case *ast.SelectorExpr:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	case *ast.CompositeLit:
		exprType := a.pass.TypesInfo.TypeOf(expr)
		msg = fmt.Sprintf("%s literal is not wrapped with stacked", typeToString(exprType, a.pass.Pkg))
	case *ast.UnaryExpr:
		switch subExpr := expr.X.(type) {
		case *ast.CompositeLit:
			exprType := a.pass.TypesInfo.TypeOf(subExpr)
			msg = fmt.Sprintf("%s literal is not wrapped with stacked", typeToString(exprType, a.pass.Pkg))
		default:
			if expr.Op == token.ARROW {
				msg = fmt.Sprintf("error received from %s is not wrapped with stacked", exprToString(expr.X))
			} else {
				msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
			}
		}
	case *ast.CallExpr:
		if a.isTypeConversion(expr) {
			msg = fmt.Sprintf("value converted to error type %s is not wrapped with stacked", exprToString(expr.Fun))
		} else {
			msg = fmt.Sprintf("error returned by %s is not wrapped with stacked", exprToString(expr.Fun))
			isIteratorPull, valueCount = a.returnValueCount(expr)
			errorReturnIndex = a.errorReturnIndex(expr)
		}
	}

	var suggestedFixes []analysis.SuggestedFix
	if autoWrap && valueCount <= 5 && (isIteratorPull || errorReturnIndex == valueCount-1) {
		wrapPull := ""
		if isIteratorPull {
			wrapPull = "Pull"
		}

		wrapValueCountString := ""
		if valueCount > 1 {
			wrapValueCountString = strconv.Itoa(valueCount)
		}

		wrapFuncName := fmt.Sprintf("stacked.Wrap%s%s", wrapPull, wrapValueCountString)

		suggestedFixes = []analysis.SuggestedFix{{
			Message: fmt.Sprintf("Wrap with %s", wrapFuncName),
			TextEdits: []analysis.TextEdit{{
				Pos:     expr.Pos(),
				End:     expr.End(),
				NewText: []byte(fmt.Sprintf("%s(%s)", wrapFuncName, exprToString(expr))),
			}},
		}}

		addMissingImport := addMissingStackedImport(a.currentFile())
		if addMissingImport != nil {
			suggestedFixes[0].TextEdits = append(suggestedFixes[0].TextEdits, *addMissingImport)
		}
	}

	a.pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func (a *analyzer) reportIterator(expr ast.Expr, valueCount int) {
	var msg string

	expr = ast.Unparen(expr)
	switch expr := expr.(type) {
	case *ast.FuncLit:
		msg = fmt.Sprintf("iterator literal is not wrapped with stacked")
	case *ast.UnaryExpr:
		if expr.Op == token.ARROW {
			msg = fmt.Sprintf("iterator received from %s is not wrapped with stacked", exprToString(expr.X))
		} else {
			msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
		}
	case *ast.CallExpr:
		if a.isTypeConversion(expr) {
			msg = fmt.Sprintf("value converted to iterator type %s is not wrapped with stacked", exprToString(expr.Fun))
		} else {
			msg = fmt.Sprintf("iterator returned by %s is not wrapped with stacked", exprToString(expr.Fun))
		}
	default:
		msg = fmt.Sprintf("%s is not wrapped with stacked", exprToString(expr))
	}

	var suggestedFixes []analysis.SuggestedFix
	if valueCount <= 2 {
		wrapValueCountString := ""
		if valueCount > 1 {
			wrapValueCountString = strconv.Itoa(valueCount)
		}

		wrapFuncName := fmt.Sprintf("stacked.WrapSeq%s", wrapValueCountString)

		suggestedFixes = []analysis.SuggestedFix{{
			Message: fmt.Sprintf("Wrap with %s", wrapFuncName),
			TextEdits: []analysis.TextEdit{{
				Pos:     expr.Pos(),
				End:     expr.End(),
				NewText: []byte(fmt.Sprintf("%s(%s)", wrapFuncName, exprToString(expr))),
			}},
		}}

		addMissingImport := addMissingStackedImport(a.currentFile())
		if addMissingImport != nil {
			suggestedFixes[0].TextEdits = append(suggestedFixes[0].TextEdits, *addMissingImport)
		}
	}

	a.pass.Report(analysis.Diagnostic{
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

func (a *analyzer) shouldWrap(expr ast.Expr) bool {
	if a.isIgnoredLine(expr) {
		return false
	}

	expr = ast.Unparen(expr)
	switch expr := expr.(type) {
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

func (a *analyzer) shouldWrapIdent(ident *ast.Ident) bool {
	obj := a.pass.TypesInfo.ObjectOf(ident)

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return implementsError(variable.Type()) && variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()
}

func (a *analyzer) shouldWrapSelector(expr *ast.SelectorExpr) bool {
	obj := a.pass.TypesInfo.ObjectOf(expr.Sel)

	variable, ok := obj.(*types.Var)
	if !ok {
		return false
	}

	return implementsError(variable.Type()) && variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()
}

func (a *analyzer) shouldWrapCompositeLit(lit *ast.CompositeLit) bool {
	return implementsError(a.pass.TypesInfo.TypeOf(lit))
}

func (a *analyzer) shouldWrapUnary(expr *ast.UnaryExpr) bool {
	isLiteralOrChannelReceive := false
	switch expr.X.(type) {
	case *ast.CompositeLit:
		isLiteralOrChannelReceive = true
	default:
		if expr.Op == token.ARROW {
			isLiteralOrChannelReceive = true
		}
	}

	return isLiteralOrChannelReceive && implementsError(a.pass.TypesInfo.TypeOf(expr))
}

func (a *analyzer) shouldWrapCall(call *ast.CallExpr) bool {
	if a.isWrapCall(call) {
		return false
	}

	if a.isTypeConversion(call) {
		return a.returnsError(call)
	}

	if isFunctionLiteral(call) {
		return a.returnsError(call)
	}

	if a.isInternalCall(call) {
		return false
	}

	if a.isIgnoredCall(call) {
		return false
	}

	return a.returnsError(call)
}

func (a *analyzer) isWrapCall(call *ast.CallExpr) bool {
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

	obj := a.pass.TypesInfo.ObjectOf(ident)

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	return pkg.Imported().Path() == stackedImportPath
}

func (a *analyzer) isTypeConversion(call *ast.CallExpr) bool {
	callType, ok := a.pass.TypesInfo.Types[call.Fun]
	return ok && callType.IsType()
}

func (a *analyzer) isInternalCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		obj := a.pass.TypesInfo.ObjectOf(fun)

		_, isFunc := obj.(*types.Func)
		if !isFunc {
			return false
		}

		return obj.Pkg() != nil
	case *ast.SelectorExpr:
		obj := a.pass.TypesInfo.ObjectOf(fun.Sel)

		sel := a.pass.TypesInfo.Selections[fun]
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

		if a.config.isPackageTreatedAsExternal(pkg.Path()) {
			return false
		}

		return strings.HasPrefix(pkg.Path(), a.pass.Module.Path)
	}

	return false
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
	if obj == nil {
		return false
	}

	if obj.Pkg() == nil || obj.Pkg().Path() != "fmt" || obj.Name() != "Errorf" {
		return false
	}

	if len(call.Args) == 0 {
		return false
	}

	formatArg := call.Args[0]
	tv, ok := a.pass.TypesInfo.Types[formatArg]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
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

	tuple, ok := returnType.(*types.Tuple)
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

func (a *analyzer) returnValueCount(call *ast.CallExpr) (bool, int) {
	callType := a.pass.TypesInfo.TypeOf(call)

	tuple, ok := callType.(*types.Tuple)
	if ok {
		if a.isIteratorPull(call) {
			return true, tuple.Len() - 1
		}

		return false, tuple.Len()
	}

	return false, 1
}

func (a *analyzer) isIteratorPull(call *ast.CallExpr) bool {
	funType := a.pass.TypesInfo.TypeOf(call.Fun)
	if funType == nil {
		return false
	}

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

func (a *analyzer) isErrorCheckCall(call *ast.CallExpr) (bool, int) {
	isErrorCheckCall, functionName := a.isCalledFunctionInSet(call, a.config.isCheckFunction)
	if !isErrorCheckCall {
		return false, 0
	}

	argumentIndex := a.config.checkFunctionArgument(functionName) - 1

	return true, argumentIndex
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
	if obj == nil {
		return false, ""
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return false, ""
	}

	pkg := fn.Pkg()
	if pkg == nil {
		return false, ""
	}

	fullName := pkg.Path()

	sig, ok := fn.Type().(*types.Signature)
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

	fullName += "." + fn.Name()

	return isFunctionInSet(fullName), fullName
}

func (a *analyzer) shouldWrapIterator(expr ast.Expr) (bool, int) {
	callExpr, isCall := expr.(*ast.CallExpr)
	if isCall {
		if a.isWrapCall(callExpr) {
			return false, 0
		}
	}

	exprType := a.pass.TypesInfo.TypeOf(expr)

	sig, ok := types.Unalias(exprType).Underlying().(*types.Signature)
	if !ok {
		return false, 0
	}

	if !(sig.Params().Len() == 1 && sig.Results().Len() == 0) {
		return false, 0
	}

	yieldSig, ok := types.Unalias(sig.Params().At(0).Type()).Underlying().(*types.Signature)
	if !ok {
		return false, 0
	}

	if yieldSig.Results().Len() != 1 {
		return false, 0
	}
	resType, ok := yieldSig.Results().At(0).Type().Underlying().(*types.Basic)
	if !ok || resType.Kind() != types.Bool {
		return false, 0
	}

	yieldParams := yieldSig.Params()
	if !(yieldParams.Len() == 1 || yieldParams.Len() == 2) {
		return false, 0
	}

	for i := 0; i < yieldParams.Len(); i++ {
		if implementsError(yieldParams.At(i).Type()) {
			return true, yieldParams.Len()
		}
	}

	return false, 0
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

func (a *analyzer) isErrorExpectedInCallArgs(call *ast.CallExpr, argIndex int) bool {
	funType := a.pass.TypesInfo.TypeOf(call.Fun)

	sig, ok := funType.Underlying().(*types.Signature)
	if ok {
		var argType types.Type
		params := sig.Params()

		if sig.Variadic() && argIndex >= params.Len()-1 {
			lastParam := params.At(params.Len() - 1).Type()
			sliceType := lastParam.(*types.Slice)
			argType = sliceType.Elem()
		} else if argIndex < params.Len() {
			argType = params.At(argIndex).Type()
		}

		return isError(argType)
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
