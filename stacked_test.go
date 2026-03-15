package stacked

import (
	"errors"
	"iter"
	"os"
	"testing"
)

func TestRecover(t *testing.T) {
	Recover(func() {
		var p *os.PathError
		panic(p)
	}, func(err error) {
		t.Log("err", err.Error(), err != nil)
	}, false)
}

func TestIteratorPull(t *testing.T) {
	err := Wrap(errors.New("error"))
	t.Log(StackTrace(err))

	seq := func(yield func(err error) bool) {
		yield(errors.New("seq error"))
	}

	for err := range WrapSeq(seq) {
		t.Log(StackTrace(err))
	}

	yield := func(err error) bool {
		t.Log(StackTrace(err))
		return false
	}
	WrapSeq(seq)(yield)

	next, stop := iter.Pull(seq)
	defer stop()

	err, _ = WrapPull(next())
	t.Log(StackTrace(err))
}
