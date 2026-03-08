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

	seq := WrapSeq(func(yield func(err error) bool) {
		yield(errors.New("seq error"))
	})

	for err := range seq {
		t.Log(StackTrace(err))
	}

	next, stop := iter.Pull(seq)
	defer stop()

	err, _ = WrapPull(next())
	t.Log(StackTrace(err))
}
