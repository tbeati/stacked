package a

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"net"
	"os"

	"github.com/tbeati/stacked"

	"testdata/b"
	"testdata/generated"
)

var errTopLevel = generated.ErrGlobal

func multiDeclarations() {
	{
		var (
			err        = os.Chdir("/") // want "^error returned by os.Chdir is not wrapped with stacked$"
			_, _       = 0, ""
			err1, err2 error
		)
		_, _, _ = err, err1, err2
	}
	{
		var (
			err        = stacked.Wrap(os.Chdir("/"))
			_, _       = 0, ""
			err1, err2 error
		)
		_, _, _ = err, err1, err2
	}

	{
		var (
			_, err     = os.Open("") // want "^error returned by os.Open is not wrapped with stacked$"
			_, _       = 0, ""
			err1, err2 error
		)
		_, _, _ = err, err1, err2
	}
	{
		var (
			_, err     = stacked.Wrap2(os.Open(""))
			_, _       = 0, ""
			err1, err2 error
		)
		_, _, _ = err, err1, err2
	}
}

func conditionDeclarations() {
	if err := os.Chdir("/"); err != nil { // want "^error returned by os.Chdir is not wrapped with stacked$"
	}
	if err := stacked.Wrap(os.Chdir("/")); err != nil {
	}

	switch err := os.Chdir("/"); err { // want "^error returned by os.Chdir is not wrapped with stacked$"
	}
	switch err := stacked.Wrap(os.Chdir("/")); err {
	}
}

func complexExpressions() {
	var err error
	_ = err

	var funcStructFieldExpr []struct {
		f func() error
	}
	err = funcStructFieldExpr[0].f() // want "^error returned by funcStructFieldExpr\\[0\\]\\.f is not wrapped with stacked$"
	err = stacked.Wrap(funcStructFieldExpr[0].f())

	var funcSliceExpr []func() error
	err = funcSliceExpr[0]() // want "^error returned by funcSliceExpr\\[0\\] is not wrapped with stacked$"
	err = stacked.Wrap(funcSliceExpr[0]())

	var chanStructFieldExpr []struct {
		c chan error
	}
	err = <-chanStructFieldExpr[0].c // want "^error received from chanStructFieldExpr\\[0\\]\\.c is not wrapped with stacked$"
	err = stacked.Wrap(<-chanStructFieldExpr[0].c)

	var chanSliceExpr []chan error
	err = <-chanSliceExpr[0] // want "^error received from chanSliceExpr\\[0\\] is not wrapped with stacked$"
	err = stacked.Wrap(<-chanSliceExpr[0])
}

func newError() {
	var err error
	_ = err

	err = new(structError) // want "^error returned by new is not wrapped with stacked$"
	err = stacked.Wrap(new(structError))

	err = new(generated.StructError) // want "^error returned by new is not wrapped with stacked$"
	err = stacked.Wrap(new(generated.StructError))
}

type doesNotImplementsError struct{}

type implementsError doesNotImplementsError

func (e implementsError) Error() string {
	return "error"
}

func typeConversionArg() {
	var err error
	_ = err

	err = error(generated.StructError{}) // want "^generated.StructError literal is not wrapped with stacked$"

	err = implementsError(doesNotImplementsError{}) // want "^value converted to error type implementsError is not wrapped with stacked$"
	err = stacked.Wrap(implementsError(doesNotImplementsError{}))

	const errMessage = "error"
	err = net.UnknownNetworkError(errMessage) // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	err = stacked.Wrap(net.UnknownNetworkError(errMessage))
}

func localConst() {
	var err net.UnknownNetworkError
	_ = err

	const errMessage = "error"
	err = errMessage // want "^errMessage is not wrapped with stacked$"
}

func externalConst() {
	var err net.UnknownNetworkError
	_ = err

	err = generated.ErrMessage // want "^generated.ErrMessage is not wrapped with stacked$"
}

func variadic(errs ...error) {
}

func variadic2(n int, errs ...error) {
}

func variadicFunctionArg() {
	variadic(errors.New("")) // want "^error returned by errors.New is not wrapped with stacked$"
	variadic(stacked.Wrap(errors.New("")))

	var f func() (int, error)
	variadic2(f()) // want "^error returned by f is not wrapped with stacked$"
	variadic2(stacked.Wrap2(f()))
}

func alternativeGenDecl() {
	const err = "error"
	type t struct{}
}

func multiErrorDeclarations() {
	var err1, err2 = errors.New("error"), errors.New("error") // want "^assignment to multiple error variables$"
	err3, err4 := errors.New("error"), errors.New("error")    // want "^assignment to multiple error variables$"
	err4, err3, err2, err1 = err1, err2, err3, err4           // want "^assignment to multiple error variables$"
	_, _, _, _ = err1, err2, err3, err4
}

func iteratorYieldsMultipleErrors() {
	for range func(yield func(e1, e2 error) bool) {} { // want "^iterator yields multiple errors$"
	}
}

func innerParenthesis() {
	var err error
	_ = err

	err = <-(errChan) // want "^error received from errChan is not wrapped with stacked$"
	err = stacked.Wrap(<-(errChan))
	err = &(fs.PathError{Path: "error"}) // want "^fs.PathError literal is not wrapped with stacked$"
	err = stacked.Wrap(&(fs.PathError{Path: "error"}))
	err = (os.Chdir)("/") // want "^error returned by os.Chdir is not wrapped with stacked$"
	err = stacked.Wrap((os.Chdir)("/"))
}

func functionLiteral() {
	var err error
	_ = err

	err = func() error {
		return err
	}()
}

func wrapNameCollision() {
	var Wrap = func(err error) error {
		return err
	}

	var err error
	_ = err

	err = Wrap(err) // want "^error returned by Wrap is not wrapped with stacked$"
}

type errorSeq func(yield func(err error) bool)

func iteratorTypeConversion() {
	var seq iter.Seq[error]

	for range errorSeq(seq) { // want "^value converted to iterator type errorSeq is not wrapped with stacked$"
	}
}

func callErrorNotLast2() (error, int) {
	var errorNotLast func() (error, int)
	return errorNotLast() // want "^error returned by errorNotLast is not wrapped with stacked$"
}

func callErrorNotLast3() (error, int, int) {
	var errorNotLast func() (error, int, int)
	return errorNotLast() // want "^error returned by errorNotLast is not wrapped with stacked$"
}

func iteratorErrorNotLast() {
	var seq iter.Seq2[error, int]
	for range seq { // want "^seq is not wrapped with stacked$"
	}
}

func iteratorPullAsArg() {
	var f func(int, error, bool)
	var pull func() (int, error, bool)

	f(pull()) // want "^error returned by pull is not wrapped with stacked$"
}

func compositeLiterals() {
	_ = structWithErrorField{
		generated.ErrGlobal, // want "^generated.ErrGlobal is not wrapped with stacked$"
	}

	_ = [1]error{
		generated.ErrGlobal, // want "^generated.ErrGlobal is not wrapped with stacked$"
	}
}

func rangeWithoutIterator() {
	for range 1 {
	}
}

func anyDestinationCallArg() {
	_ = fmt.Errorf("error: %w", os.Chdir("/")) // want "^error returned by os.Chdir is not wrapped with stacked$"
	_ = fmt.Errorf("error: %w", stacked.Wrap(os.Chdir("/")))
}

func anyDestinationAssignment() {
	var x any
	x = os.Chdir("/") // want "^error returned by os.Chdir is not wrapped with stacked$"
	x = stacked.Wrap(os.Chdir("/"))
	_ = x
}

func anyDestinationReturn() any {
	return os.Chdir("/") // want "^error returned by os.Chdir is not wrapped with stacked$"
	return stacked.Wrap(os.Chdir("/"))
}

func anyDestinationLiteral() {
	_ = []any{os.Chdir("/")} // want "^error returned by os.Chdir is not wrapped with stacked$"
	_ = []any{stacked.Wrap(os.Chdir("/"))}
}

func genericCalls() {
	var err error
	_ = err

	err = generated.GenericOneParam[int](0) // want "^error returned by generated.GenericOneParam\\[int\\] is not wrapped with stacked$"
	err = stacked.Wrap(generated.GenericOneParam[int](0))

	err = generated.GenericTwoParams[int, string](0, "") // want "^error returned by generated.GenericTwoParams\\[int, string\\] is not wrapped with stacked$"
	err = stacked.Wrap(generated.GenericTwoParams[int, string](0, ""))

	err = b.GenericOneParam[int](0)
	err = b.GenericTwoParams[int, string](0, "")

	_, err = stacked.Wrap2[[]byte](os.ReadFile(""))
}
