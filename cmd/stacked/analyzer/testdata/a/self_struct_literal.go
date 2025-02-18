package a

import (
	"github.com/beati/stacked"
)

func structLiteralAssignmentSelf() {
	var err error

	err = structError{message: "error"} // want "structError literal is not wrapped with stacked"

	err = stacked.Wrap(structError{message: "error"})

	err = structError{message: "error"}
	err = stacked.Wrap(err)

	_, err = 0, structError{message: "error"} // want "structError literal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(structError{message: "error"})

	_, err = 0, structError{message: "error"}
	err = stacked.Wrap(err)
}

func structLiteralDeclarationSelf() {
	{
		var err = structError{message: "error"} // want "structError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(structError{message: "error"})
		_ = err
	}

	{
		var err error = structError{message: "error"} // want "structError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err error = structError{message: "error"}
		err = stacked.Wrap(err)
	}

	{
		var err = error(structError{message: "error"}) // want "structError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, structError{message: "error"} // want "structError literal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(structError{message: "error"})
		_ = err
	}

	{
		var _, err = 0, error(structError{message: "error"}) // want "structError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func structLiteralShortDeclarationSelf() {
	{
		err := structError{message: "error"} // want "structError literal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(structError{message: "error"})
		_ = err
	}

	{
		err := error(structError{message: "error"}) // want "structError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, structError{message: "error"} // want "structError literal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(structError{message: "error"})
		_ = err
	}

	{
		_, err := 0, error(structError{message: "error"}) // want "structError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func structLiteralReturnSingleSelf() error {
	return structError{message: "error"} // want "structError literal is not wrapped with stacked"
	return stacked.Wrap(structError{message: "error"})
}

func structLiteralReturnMultipleSelf() (int, error) {
	return 0, structError{message: "error"} // want "structError literal is not wrapped with stacked"
	return 0, stacked.Wrap(structError{message: "error"})
}

func structLiteralArgumentSelf() {
	functionWithIntErrorArgument(0, structError{message: "error"}) // want "structError literal is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(structError{message: "error"}))
}

func structLiteralCompositeLiteralSelf() {
	_ = structWithErrorField{
		err: structError{message: "error"}, // want "structError literal is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(structError{message: "error"}),
	}

	_ = []error{structError{message: "error"}} // want "structError literal is not wrapped with stacked"
	_ = []error{stacked.Wrap(structError{message: "error"})}

	_ = map[string]error{"": structError{message: "error"}} // want "structError literal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(structError{message: "error"})}
}
