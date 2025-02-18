package a

import (
	"io/fs"

	"github.com/tbeati/stacked"
)

func globalVariableAssignmentExternal() {
	var err error

	err = fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"

	err = stacked.Wrap(fs.ErrNotExist)

	err = fs.ErrNotExist
	err = stacked.Wrap(err)

	_, err = 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"

	_, err = 0, stacked.Wrap(fs.ErrNotExist)

	_, err = 0, fs.ErrNotExist
	err = stacked.Wrap(err)
}

func globalVariableDeclarationExternal() {
	{
		var err = fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(fs.ErrNotExist)
		_ = err
	}

	{
		var err = fs.ErrNotExist
		err = stacked.Wrap(err)
	}

	{
		var _, err = 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(fs.ErrNotExist)
		_ = err
	}

	{
		var _, err = 0, fs.ErrNotExist
		err = stacked.Wrap(err)
	}
}

func globalVariableShortDeclarationExternal() {
	{
		err := fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(fs.ErrNotExist)
		_ = err
	}

	{
		err := fs.ErrNotExist
		err = stacked.Wrap(err)
	}

	{
		_, err := 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(fs.ErrNotExist)
		_ = err
	}

	{
		_, err := 0, fs.ErrNotExist
		err = stacked.Wrap(err)
	}
}

func globalVariableReturnSingleExternal() error {
	return fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	return stacked.Wrap(fs.ErrNotExist)
}

func globalVariableReturnMultipleExternal() (int, error) {
	return 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	return 0, stacked.Wrap(fs.ErrNotExist)
}

func globalVariableArgumentExternal() {
	functionWithIntErrorArgument(0, fs.ErrNotExist) // want "fs.ErrNotExist is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(fs.ErrNotExist))
}

func globalVariableCompositeLiteralExternal() {
	_ = structWithErrorField{
		err: fs.ErrNotExist, // want "fs.ErrNotExist is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(fs.ErrNotExist),
	}

	_ = []error{fs.ErrNotExist} // want "fs.ErrNotExist is not wrapped with stacked"
	_ = []error{stacked.Wrap(fs.ErrNotExist)}

	_ = map[string]error{"": fs.ErrNotExist} // want "fs.ErrNotExist is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(fs.ErrNotExist)}
}
