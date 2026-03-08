package generated

import (
	"errors"
	"iter"
)

var ErrGlobal = errors.New("error")

type StringError string

func (err StringError) Error() string {
	return string(err)
}

var ErrGlobalConcreteType = StringError("error")
var ErrGlobalConcreteTypePointer = &ErrGlobalConcreteType

type StructError struct {
	Message string
}

func (err StructError) Error() string {
	return err.Message
}

func SingleReturn() error {
	return nil
}

func MultipleReturn() (int, error) {
	return 0, nil
}

type StructWithMethods struct{}

func (s *StructWithMethods) SingleReturn() error {
	return nil
}

func (s *StructWithMethods) MultipleReturn() (int, error) {
	return 0, nil
}

type Interface interface {
	SingleReturn() error
	MultipleReturn() (int, error)
}

func ReturnConcreteType() StringError {
	return "error"
}

func ReturnConcreteTypePointer() *StringError {
	var err = StringError("error")
	return &err
}

func IgnoredFunction(err error) error {
	return err
}

type IgnoredStruct struct{}

func (s *IgnoredStruct) IgnoredMethod(err error) error {
	return err
}

type IgnoredInterface interface {
	IgnoredMethod() error
}

func Seq(yield func(err error) bool) {}

func Seq2(yield func(n int, err error) bool) {}

func Iterator() iter.Seq[error] {
	return Seq
}

func Iterator2() iter.Seq2[int, error] {
	return Seq2
}
