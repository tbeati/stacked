package a

import (
	"github.com/tbeati/stacked"
)

func globalVariableAssignmentSelf() {
	var err error

	err = errGlobal // want "errGlobal is not wrapped with stacked"

	err = stacked.Wrap(errGlobal)

	err = errGlobal
	err = stacked.Wrap(err)

	_, err = 0, errGlobal // want "errGlobal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(errGlobal)

	_, err = 0, errGlobal
	err = stacked.Wrap(err)
}

func globalVariableDeclarationSelf() {
	{
		var err = errGlobal // want "errGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(errGlobal)
		_ = err
	}

	{
		var err = errGlobal
		err = stacked.Wrap(err)
	}

	{
		var _, err = 0, errGlobal // want "errGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(errGlobal)
		_ = err
	}

	{
		var _, err = 0, errGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableShortDeclarationSelf() {
	{
		err := errGlobal // want "errGlobal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(errGlobal)
		_ = err
	}

	{
		err := errGlobal
		err = stacked.Wrap(err)
	}

	{
		_, err := 0, errGlobal // want "errGlobal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(errGlobal)
		_ = err
	}

	{
		_, err := 0, errGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableReturnSingleSelf() error {
	return errGlobal // want "errGlobal is not wrapped with stacked"
	return stacked.Wrap(errGlobal)
}

func globalVariableReturnMultipleSelf() (int, error) {
	return 0, errGlobal // want "errGlobal is not wrapped with stacked"
	return 0, stacked.Wrap(errGlobal)
}

func globalVariableArgumentSelf() {
	functionWithIntErrorArgument(0, errGlobal) // want "errGlobal is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(errGlobal))
}

func globalVariableCompositeLiteralSelf() {
	_ = structWithErrorField{
		err: errGlobal, // want "errGlobal is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(errGlobal),
	}

	_ = []error{errGlobal} // want "errGlobal is not wrapped with stacked"
	_ = []error{stacked.Wrap(errGlobal)}

	_ = map[string]error{"": errGlobal} // want "errGlobal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(errGlobal)}
}
