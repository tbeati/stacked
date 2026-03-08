package a

import (
	"github.com/tbeati/stacked"

	"testdata/b"
)

func assignmentInternal() {
	var err error
	_ = err

	err = b.SingleReturn()
	_, err = b.MultipleReturn()

	localFuncSingle := b.SingleReturn
	err = localFuncSingle() // want "error returned by localFuncSingle is not wrapped with stacked"
	err = stacked.Wrap(localFuncSingle())

	localFuncMultiple := b.MultipleReturn
	_, err = localFuncMultiple() // want "error returned by localFuncMultiple is not wrapped with stacked"
	_, err = stacked.Wrap2(localFuncMultiple())

	s := b.StructWithMethods{}
	err = s.SingleReturn()
	_, err = s.MultipleReturn()

	var i b.Interface
	err = i.SingleReturn() // want "error returned by i.SingleReturn is not wrapped with stacked"
	err = stacked.Wrap(i.SingleReturn())
	_, err = i.MultipleReturn() // want "error returned by i.MultipleReturn is not wrapped with stacked"
	_, err = stacked.Wrap2(i.MultipleReturn())

	err = b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
	err = stacked.Wrap(b.ErrGlobal)

	err = b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
	err = stacked.Wrap(b.StringError("error"))

	err = b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
	err = stacked.Wrap(b.StructError{Message: "error"})
}
