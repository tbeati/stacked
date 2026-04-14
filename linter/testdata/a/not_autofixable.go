package a

import (
	"testdata/generated"
)

func functionWithConcreteStringErrorArgument(err generated.StringError)                 {}
func functionWithIntConcreteStringErrorArgument(n int, err generated.StringError)       {}
func functionWithStringConcreteStringErrorArgument(s string, err generated.StringError) {}

type structWithConcreteStringErrorField struct {
	err generated.StringError
}

func notAutoFixableFunctionCallAssignment() {
	var err generated.StringError
	_ = err

	err = generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"

	_, err = 0, generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableFunctionCallDeclaration() {
	{
		var err = generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
		_ = err
	}
	{
		var _, err = 0, generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
		_ = err
	}
}

func notAutoFixableFunctionCallShortDeclaration() {
	err := generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
	_ = err
}

func notAutoFixableFunctionCallReturnSingle() generated.StringError {
	return generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableFunctionCallReturnMultiple() (int, generated.StringError) {
	return 0, generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableFunctionCallArgument() {
	functionWithConcreteStringErrorArgument(generated.ReturnConcreteType()) // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"

	functionWithStringConcreteStringErrorArgument("", generated.ReturnConcreteType()) // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"

	functionWithIntConcreteStringErrorArgument(0, generated.ReturnConcreteType()) // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableFunctionCallCompositeLiteral() {
	_ = structWithConcreteStringErrorField{
		err: generated.ReturnConcreteType(), // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
	}

	_ = []generated.StringError{generated.ReturnConcreteType()} // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"

	_ = map[string]generated.StringError{"": generated.ReturnConcreteType()} // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableFunctionCallChannelSend() {
	var errChan chan generated.StringError

	errChan <- generated.ReturnConcreteType() // want "^error returned by generated.ReturnConcreteType is not wrapped with stacked$"
}

func notAutoFixableIteratorRange() {
	for range func(yield func(err stringError) bool) {} { // want "^iterator literal is not wrapped with stacked$"
	}

	for range func(yield func(n int, err stringError) bool) {} { // want "^iterator literal is not wrapped with stacked$"
	}
}
