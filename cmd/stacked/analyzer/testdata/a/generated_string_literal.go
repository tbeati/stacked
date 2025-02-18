package a

import (
	"github.com/tbeati/stacked"

	"testdata/generated"
)

func stringLiteralAssignmentGenerated() {
	var err error

	err = generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"

	err = stacked.Wrap(generated.StringError("error"))

	err = generated.StringError("error")
	err = stacked.Wrap(err)

	_, err = 0, generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"

	_, err = 0, stacked.Wrap(generated.StringError("error"))

	_, err = 0, generated.StringError("error")
	err = stacked.Wrap(err)
}

func stringLiteralDeclarationGenerated() {
	{
		var err = generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(generated.StringError("error"))
		_ = err
	}

	{
		var err error = generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var err error = generated.StringError("error")
		err = stacked.Wrap(err)
	}

	{
		var err = error(generated.StringError("error")) // want "value converted to error type generated.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(generated.StringError("error"))
		_ = err
	}

	{
		var _, err = 0, error(generated.StringError("error")) // want "value converted to error type generated.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralShortDeclarationGenerated() {
	{
		err := generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(generated.StringError("error"))
		_ = err
	}

	{
		err := error(generated.StringError("error")) // want "value converted to error type generated.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(generated.StringError("error"))
		_ = err
	}

	{
		_, err := 0, error(generated.StringError("error")) // want "value converted to error type generated.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralReturnSingleGenerated() error {
	return generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
	return stacked.Wrap(generated.StringError("error"))
}

func stringLiteralReturnMultipleGenerated() (int, error) {
	return 0, generated.StringError("error") // want "value converted to error type generated.StringError is not wrapped with stacked"
	return 0, stacked.Wrap(generated.StringError("error"))
}

func stringLiteralArgumentGenerated() {
	functionWithIntErrorArgument(0, generated.StringError("error")) // want "value converted to error type generated.StringError is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(generated.StringError("error")))
}

func stringLiteralCompositeLiteralGenerated() {
	_ = structWithErrorField{
		err: generated.StringError("error"), // want "value converted to error type generated.StringError is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(generated.StringError("error")),
	}

	_ = []error{generated.StringError("error")} // want "value converted to error type generated.StringError is not wrapped with stacked"
	_ = []error{stacked.Wrap(generated.StringError("error"))}

	_ = map[string]error{"": generated.StringError("error")} // want "value converted to error type generated.StringError is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(generated.StringError("error"))}
}
