package a

import (
	"io/fs"

	"github.com/tbeati/stacked"
)

func globalVariableAssignmentExternal() {
	var err error
	_ = err

	err = fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	err = stacked.Wrap(fs.ErrNotExist)

	_, err = 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	_, err = 0, stacked.Wrap(fs.ErrNotExist)
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
		var _, err = 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(fs.ErrNotExist)
		_ = err
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
		_, err := 0, fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(fs.ErrNotExist)
		_ = err
	}
}

func globalVariableReturn1External() error {
	return fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	return stacked.Wrap(fs.ErrNotExist)
}

func globalVariableReturn2External() (int, error) {
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

func globalVariableChannelSendExternal() {
	var errChan chan error

	errChan <- fs.ErrNotExist // want "fs.ErrNotExist is not wrapped with stacked"
	errChan <- stacked.Wrap(fs.ErrNotExist)
}
