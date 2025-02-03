package a

import (
	"github.com/beati/stacked"

	"testdata/b"
)

func callInternalFuncAssignment() error {
	err := b.F()
	if err != nil {
		return err
	}

	err = b.F()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	err = stacked.Wrap(b.F())
	if err != nil {
		return err
	}

	return nil
}

func callInternalMethodAssignment() error {
	s := b.S{}
	err := s.F()
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
