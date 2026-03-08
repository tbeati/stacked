package a

import (
	"errors"
	"fmt"
	"io/fs"

	"testdata/generated"
)

// errors As AsType Is Unwrap Join

func assignmentIgnoredFunction() {
	var err error
	_ = err

	err = generated.IgnoredFunction(err)

	s := generated.IgnoredStruct{}
	err = s.IgnoredMethod(err)
	err = (&s).IgnoredMethod(err)

	var i generated.IgnoredInterface
	err = i.IgnoredMethod()

	err = fmt.Errorf("wrapping %w", err)
	err = fmt.Errorf("not wrapping") // want "error returned by fmt.Errorf is not wrapped with stacked"

	err = errors.Join(err)
	err = errors.Unwrap(err)
}

func errorCheckFunctions() {
	var err error
	errors.Is(err, fs.ErrNotExist)
	errors.Is(err, (fs.ErrNotExist))
	errors.As(err, &fs.PathError{})
}
