package stacked

import (
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
