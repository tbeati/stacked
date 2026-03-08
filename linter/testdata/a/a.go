package a

import (
	"errors"
	"io/fs"
	"net/netip"
	"os"
)

type structWithErrorField struct {
	err error
}

func functionWithErrorArgument(err error)                                                           {}
func functionWithIntErrorArgument(n int, err error)                                                 {}
func functionWithStringErrorArgument(s string, err error)                                           {}
func functionWithFileErrorArgument(f *os.File, err error)                                           {}
func functionWithFileFileErrorArgument(r *os.File, w *os.File, err error)                           {}
func functionWithIntIntIntAddrPortErrorArgument(n, oobn, flags int, addr netip.AddrPort, err error) {}
func functionWithFileInfoErrorArgument(f fs.FileInfo, err error)                                    {}

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

type localInterface interface {
	SingleReturn() error
	MultipleReturn() (int, error)
}

var errChan chan error
