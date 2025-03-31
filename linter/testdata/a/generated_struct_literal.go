package a

import (
	"github.com/tbeati/stacked"

	"testdata/generated"
)

func structLiteralAssignmentGenerated() {
	var err error

	err = generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"

	err = stacked.Wrap(generated.StructError{Message: "error"})

	err = generated.StructError{Message: "error"}
	err = stacked.Wrap(err)

	_, err = 0, generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(generated.StructError{Message: "error"})

	_, err = 0, generated.StructError{Message: "error"}
	err = stacked.Wrap(err)
}

func structLiteralDeclarationGenerated() {
	{
		var err = generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(generated.StructError{Message: "error"})
		_ = err
	}

	{
		var err error = generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err error = generated.StructError{Message: "error"}
		err = stacked.Wrap(err)
	}

	{
		var err = error(generated.StructError{Message: "error"}) // want "generated.StructError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(generated.StructError{Message: "error"})
		_ = err
	}

	{
		var _, err = 0, error(generated.StructError{Message: "error"}) // want "generated.StructError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func structLiteralShortDeclarationGenerated() {
	{
		err := generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(generated.StructError{Message: "error"})
		_ = err
	}

	{
		err := error(generated.StructError{Message: "error"}) // want "generated.StructError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(generated.StructError{Message: "error"})
		_ = err
	}

	{
		_, err := 0, error(generated.StructError{Message: "error"}) // want "generated.StructError literal is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func structLiteralReturnSingleGenerated() error {
	return generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
	return stacked.Wrap(generated.StructError{Message: "error"})
}

func structLiteralReturnMultipleGenerated() (int, error) {
	return 0, generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
	return 0, stacked.Wrap(generated.StructError{Message: "error"})
}

func structLiteralArgumentGenerated() {
	functionWithIntErrorArgument(0, generated.StructError{Message: "error"}) // want "generated.StructError literal is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(generated.StructError{Message: "error"}))
}

func structLiteralCompositeLiteralGenerated() {
	_ = structWithErrorField{
		err: generated.StructError{Message: "error"}, // want "generated.StructError literal is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(generated.StructError{Message: "error"}),
	}

	_ = []error{generated.StructError{Message: "error"}} // want "generated.StructError literal is not wrapped with stacked"
	_ = []error{stacked.Wrap(generated.StructError{Message: "error"})}

	_ = map[string]error{"": generated.StructError{Message: "error"}} // want "generated.StructError literal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(generated.StructError{Message: "error"})}
}

func structLiteralReturnSinglePointerGenerated() error {
	return &generated.StructError{Message: "error"} // want "generated.StructError literal is not wrapped with stacked"
	return stacked.Wrap(&generated.StructError{Message: "error"})
}
