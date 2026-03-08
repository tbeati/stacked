package a

import (
	"iter"

	"github.com/tbeati/stacked"

	"testdata/generated"
)

func assignmentGenerated() {
	var err error
	_ = err

	err = generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
	err = stacked.Wrap(generated.SingleReturn())

	_, err = generated.MultipleReturn() // want "error returned by generated.MultipleReturn is not wrapped with stacked"
	_, err = stacked.Wrap2(generated.MultipleReturn())

	localFuncSingle := generated.SingleReturn
	err = localFuncSingle() // want "error returned by localFuncSingle is not wrapped with stacked"
	err = stacked.Wrap(localFuncSingle())

	localFuncMultiple := generated.MultipleReturn
	_, err = localFuncMultiple() // want "error returned by localFuncMultiple is not wrapped with stacked"
	_, err = stacked.Wrap2(localFuncMultiple())

	s := generated.StructWithMethods{}
	err = s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
	err = stacked.Wrap(s.SingleReturn())

	_, err = s.MultipleReturn() // want "error returned by s.MultipleReturn is not wrapped with stacked"
	_, err = stacked.Wrap2(s.MultipleReturn())

	var i generated.Interface
	err = i.SingleReturn() // want "error returned by i.SingleReturn is not wrapped with stacked"
	err = stacked.Wrap(i.SingleReturn())
	_, err = i.MultipleReturn() // want "error returned by i.MultipleReturn is not wrapped with stacked"
	_, err = stacked.Wrap2(i.MultipleReturn())

	err = generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
	err = stacked.Wrap(generated.ErrGlobal)

	err = generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
	err = stacked.Wrap(generated.StringError("error"))

	err = generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
	err = stacked.Wrap(generated.StructError{Message: "error"})

	err = generated.ReturnConcreteType() // want "error returned by generated.ReturnConcreteType is not wrapped with stacked"
	err = stacked.Wrap(generated.ReturnConcreteType())
	err = generated.ReturnConcreteTypePointer() // want "error returned by generated.ReturnConcreteTypePointer is not wrapped with stacked"
	err = stacked.Wrap(generated.ReturnConcreteTypePointer())

	err = &generated.ErrGlobalConcreteType // want "generated.ErrGlobalConcreteType is not wrapped with stacked"
	err = stacked.Wrap(generated.ErrGlobalConcreteType)
	err = generated.ErrGlobalConcreteTypePointer // want "generated.ErrGlobalConcreteTypePointer is not wrapped with stacked"
	err = stacked.Wrap(generated.ErrGlobalConcreteTypePointer)
}

func iteratorGenerated() {
	var seq iter.Seq[error]
	_ = seq

	seq = generated.Seq

	for range seq { // want "seq is not wrapped with stacked"
	}
	for range stacked.WrapSeq(seq) {
	}

	for range generated.Seq { // want "generated.Seq is not wrapped with stacked"
	}
	for range stacked.WrapSeq(generated.Seq) {
	}

	for range generated.Iterator() { // want "iterator returned by generated.Iterator is not wrapped with stacked"
	}
	for range stacked.WrapSeq(generated.Iterator()) {
	}

	for range func(yield func(err error) bool) {} { // want "iterator literal is not wrapped with stacked"
	}
	for range stacked.WrapSeq(func(yield func(err error) bool) {}) {
	}

	yield := func(err error) bool { return false }

	seq(yield) // want "seq is not wrapped with stacked"
	stacked.WrapSeq(seq)(yield)

	generated.Seq(yield) // want "generated.Seq is not wrapped with stacked"
	stacked.WrapSeq(generated.Seq)(yield)

	generated.Iterator()(yield) // want "iterator returned by generated.Iterator is not wrapped with stacked"
	stacked.WrapSeq(generated.Iterator())(yield)

	func(yield func(err error) bool) {}(yield) // want "iterator literal is not wrapped with stacked"
	stacked.WrapSeq(func(yield func(err error) bool) {})(yield)

	var seq2 iter.Seq2[int, error]
	_ = seq2

	seq2 = generated.Seq2

	for range seq2 { // want "seq2 is not wrapped with stacked"
	}
	for range stacked.WrapSeq2(seq2) {
	}

	for range generated.Seq2 { // want "generated.Seq2 is not wrapped with stacked"
	}
	for range stacked.WrapSeq2(generated.Seq2) {
	}

	for range generated.Iterator2() { // want "iterator returned by generated.Iterator2 is not wrapped with stacked"
	}
	for range stacked.WrapSeq2(generated.Iterator2()) {
	}

	for range func(yield func(n int, err error) bool) {} { // want "iterator literal is not wrapped with stacked"
	}
	for range stacked.WrapSeq2(func(yield func(n int, err error) bool) {}) {
	}

	yield2 := func(n int, err error) bool { return false }

	seq2(yield2) // want "seq2 is not wrapped with stacked"
	stacked.WrapSeq2(seq2)(yield2)

	generated.Seq2(yield2) // want "generated.Seq2 is not wrapped with stacked"
	stacked.WrapSeq2(generated.Seq2)(yield2)

	generated.Iterator2()(yield2) // want "iterator returned by generated.Iterator2 is not wrapped with stacked"
	stacked.WrapSeq2(generated.Iterator2())(yield2)

	func(yield func(n int, err error) bool) {}(yield2) // want "iterator literal is not wrapped with stacked"
	stacked.WrapSeq2(func(yield func(n int, err error) bool) {})(yield2)

	// TODO: iterator as method
}

func iteratorPullGenerated() {
	var err error
	_ = err

	var seq iter.Seq[error]
	next, _ := iter.Pull(seq)
	err, _ = next() // want "error returned by next is not wrapped with stacked"
	err, _ = stacked.WrapPull(next())

	var seq2 iter.Seq2[int, error]
	_ = seq2
	next2, _ := iter.Pull2(seq2)
	_, err, _ = next2() // want "error returned by next2 is not wrapped with stacked"
	_, err, _ = stacked.WrapPull2(next2())
}
