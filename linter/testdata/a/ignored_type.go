package a

import (
	"testdata/generated"
)

func ignoredTypeReturn() error {
	return generated.WrappedError{}
	return &generated.WrappedError{}
}

func ignoredTypeAssignment() {
	var err error
	err = generated.WrappedError{}
	err = &generated.WrappedError{}
	_ = err
}

func makeWrapped() *generated.WrappedError {
	return &generated.WrappedError{}
}

func ignoredTypeCallResult() error {
	return makeWrapped()
}

func ignoredTypeVariable() error {
	e := makeWrapped()
	return e
}

func ignoredTypeInCompositeLit() {
	_ = structWithErrorField{
		err: &generated.WrappedError{},
	}
	_ = []error{generated.WrappedError{}}
	_ = map[string]error{"": &generated.WrappedError{}}
	_ = []*generated.WrappedError{{}}
}

func ignoredGenericType() {
	var err error
	_ = err

	err = generated.GenericWrappedError[int]{}
	err = &generated.GenericWrappedError[int]{}
	err = generated.GenericWrappedError[string]{}
}
