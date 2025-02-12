package a

import (
	"github.com/beati/stacked"

	"testdata/b"
)

func globalVariableAssignmentInternal() {
	var err error

	err = b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"

	err = stacked.Wrap(b.ErrGlobal)

	err = b.ErrGlobal
	err = stacked.Wrap(err)

	_, err = 0, b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(b.ErrGlobal)

	_, err = 0, b.ErrGlobal
	err = stacked.Wrap(err)
}

func globalVariableDeclarationInternal() {
	{
		var err = b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(b.ErrGlobal)
		_ = err
	}

	{
		var err = b.ErrGlobal
		err = stacked.Wrap(err)
	}

	{
		var _, err = 0, b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(b.ErrGlobal)
		_ = err
	}

	{
		var _, err = 0, b.ErrGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableShortDeclarationInternal() {
	{
		err := b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(b.ErrGlobal)
		_ = err
	}

	{
		err := b.ErrGlobal
		err = stacked.Wrap(err)
	}

	{
		_, err := 0, b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(b.ErrGlobal)
		_ = err
	}

	{
		_, err := 0, b.ErrGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableReturnSingleInternal() error {
	return b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
	return stacked.Wrap(b.ErrGlobal)
}

func globalVariableReturnMultipleInternal() (int, error) {
	return 0, b.ErrGlobal // want "b.ErrGlobal is not wrapped with stacked"
	return 0, stacked.Wrap(b.ErrGlobal)
}

func globalVariableArgumentInternal() {
	errArgument(0, b.ErrGlobal) // want "b.ErrGlobal is not wrapped with stacked"
	errArgument(0, stacked.Wrap(b.ErrGlobal))
}

func globalVariableCompositeLiteralInternal() {
	_ = errStruct{
		err: b.ErrGlobal, // want "b.ErrGlobal is not wrapped with stacked"
	}
	_ = errStruct{
		err: stacked.Wrap(b.ErrGlobal),
	}

	_ = []error{b.ErrGlobal} // want "b.ErrGlobal is not wrapped with stacked"
	_ = []error{stacked.Wrap(b.ErrGlobal)}

	_ = map[string]error{"": b.ErrGlobal} // want "b.ErrGlobal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(b.ErrGlobal)}
}
