package a

import (
	"github.com/tbeati/stacked"
)

func assignmentSelf() {
	var err error
	_ = err

	err = singleReturn()
	_, err = multipleReturn()

	localFuncSingle := singleReturn
	err = localFuncSingle() // want "error returned by localFuncSingle is not wrapped with stacked"
	err = stacked.Wrap(localFuncSingle())

	localFuncMultiple := multipleReturn
	_, err = localFuncMultiple() // want "error returned by localFuncMultiple is not wrapped with stacked"
	_, err = stacked.Wrap2(localFuncMultiple())

	s := structWithMethods{}
	err = s.singleReturn()
	_, err = s.multipleReturn()

	var i localInterface
	err = i.SingleReturn() // want "error returned by i.SingleReturn is not wrapped with stacked"
	err = stacked.Wrap(i.SingleReturn())
	_, err = i.MultipleReturn() // want "error returned by i.MultipleReturn is not wrapped with stacked"
	_, err = stacked.Wrap2(i.MultipleReturn())

	err = errGlobal // want "errGlobal is not wrapped with stacked"
	err = stacked.Wrap(errGlobal)

	err = stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
	err = stacked.Wrap(stringError("error"))

	err = structError{message: "error"} // want "structError literal is not wrapped with stacked"
	err = stacked.Wrap(structError{message: "error"})
}
