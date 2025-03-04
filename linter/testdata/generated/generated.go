package generated

import (
	"errors"
)

var ErrGlobal = errors.New("error")

type StringError string

func (err StringError) Error() string {
	return string(err)
}

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

func IgnoredFunction(err error) error {
	return err
}

type IgnoredStruct struct{}

func (s *IgnoredStruct) IgnoredMethod(err error) error {
	return err
}
