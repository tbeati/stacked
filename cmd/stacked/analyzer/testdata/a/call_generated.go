package a

import (
	"github.com/beati/stacked"

	"testdata/generated"
)

func callGeneratedCodeFuncAssignment() error {
	err := generated.F() // want "error returned by generated.F is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = generated.F()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	err = stacked.Wrap(generated.F())
	if err != nil {
		return err
	}

	return nil
}

func callGeneratedCodeMethodAssignment() error {
	s := generated.S{}
	err := s.F() // want "error returned by s.F is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = s.F()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	err = stacked.Wrap(s.F())
	if err != nil {
		return err
	}

	return nil
}
