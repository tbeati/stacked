package a

import (
	"github.com/beati/stacked"

	"testdata/generated"
)

func globalVariableAssignmentGenerated() {
	var err error

	err = generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"

	err = stacked.Wrap(generated.ErrGlobal)

	err = generated.ErrGlobal
	err = stacked.Wrap(err)

	_, err = 0, generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(generated.ErrGlobal)

	_, err = 0, generated.ErrGlobal
	err = stacked.Wrap(err)
}

func globalVariableDeclarationGenerated() {
	{
		var err = generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(generated.ErrGlobal)
		_ = err
	}

	{
		var err = generated.ErrGlobal
		err = stacked.Wrap(err)
	}

	{
		var _, err = 0, generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(generated.ErrGlobal)
		_ = err
	}

	{
		var _, err = 0, generated.ErrGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableShortDeclarationGenerated() {
	{
		err := generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(generated.ErrGlobal)
		_ = err
	}

	{
		err := generated.ErrGlobal
		err = stacked.Wrap(err)
	}

	{
		_, err := 0, generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(generated.ErrGlobal)
		_ = err
	}

	{
		_, err := 0, generated.ErrGlobal
		err = stacked.Wrap(err)
	}
}

func globalVariableReturnSingleGenerated() error {
	return generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
	return stacked.Wrap(generated.ErrGlobal)
}

func globalVariableReturnMultipleGenerated() (int, error) {
	return 0, generated.ErrGlobal // want "generated.ErrGlobal is not wrapped with stacked"
	return 0, stacked.Wrap(generated.ErrGlobal)
}

func globalVariableArgumentGenerated() {
	functionWithIntErrorArgument(0, generated.ErrGlobal) // want "generated.ErrGlobal is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(generated.ErrGlobal))
}

func globalVariableCompositeLiteralGenerated() {
	_ = structWithErrorField{
		err: generated.ErrGlobal, // want "generated.ErrGlobal is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(generated.ErrGlobal),
	}

	_ = []error{generated.ErrGlobal} // want "generated.ErrGlobal is not wrapped with stacked"
	_ = []error{stacked.Wrap(generated.ErrGlobal)}

	_ = map[string]error{"": generated.ErrGlobal} // want "generated.ErrGlobal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(generated.ErrGlobal)}
}
