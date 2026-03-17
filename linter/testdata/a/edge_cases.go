package a

import (
	"os"

	"github.com/tbeati/stacked"
	"testdata/generated"
)

func multiDeclarations() {
	{
		var (
			err        = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
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
			_, err     = os.Open("") // want "error returned by os.Open is not wrapped with stacked"
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
	if err := os.Chdir("/"); err != nil { // want "error returned by os.Chdir is not wrapped with stacked"
	}
	if err := stacked.Wrap(os.Chdir("/")); err != nil {
	}

	switch err := os.Chdir("/"); err { // want "error returned by os.Chdir is not wrapped with stacked"
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
	err = funcStructFieldExpr[0].f() // want "error returned by funcStructFieldExpr\\[0\\]\\.f is not wrapped with stacked"
	err = stacked.Wrap(funcStructFieldExpr[0].f())

	funcSliceExpr := []func() error{}
	err = funcSliceExpr[0]() // want "error returned by funcSliceExpr\\[0\\] is not wrapped with stacked"
	err = stacked.Wrap(funcSliceExpr[0]())

	var chanStructFieldExpr []struct {
		c chan error
	}
	err = <-chanStructFieldExpr[0].c // want "error received from chanStructFieldExpr\\[0\\]\\.c is not wrapped with stacked"
	err = stacked.Wrap(<-chanStructFieldExpr[0].c)

	chanSliceExpr := []chan error{}
	err = <-chanSliceExpr[0] // want "error received from chanSliceExpr\\[0\\] is not wrapped with stacked"
	err = stacked.Wrap(<-chanSliceExpr[0])
}

func newError() {
	var err error
	_ = err

	err = new(generated.StructError) // want "error returned by new is not wrapped with stacked"
}
