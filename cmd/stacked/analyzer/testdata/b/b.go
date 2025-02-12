package b

import (
	"errors"
)

var ErrGlobal = errors.New("b error")

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

type S struct{}

func (s *S) SingleReturn() error {
	return nil
}

func (s *S) MultipleReturn() (int, error) {
	return 0, nil
}
