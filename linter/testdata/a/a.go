package a

import "errors"

type structWithErrorField struct {
	err error
}

func functionWithIntErrorArgument(n int, err error)       {}
func functionWithStringErrorArgument(s string, err error) {}

var errGlobal = errors.New("error")

type stringError string

func (err stringError) Error() string {
	return string(err)
}

type structError struct {
	message string
}

func (err structError) Error() string {
	return err.message
}

func singleReturn() error {
	return nil
}

func multipleReturn() (int, error) {
	return 0, nil
}

type structWithMethods struct{}

func (s *structWithMethods) singleReturn() error {
	return nil
}

func (s *structWithMethods) multipleReturn() (int, error) {
	return 0, nil
}
