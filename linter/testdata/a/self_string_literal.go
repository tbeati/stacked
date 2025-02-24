package a

import (
	"github.com/tbeati/stacked"
)

func stringLiteralAssignmentSelf() {
	var err error

	err = stringError("error") // want "value converted to error type stringError is not wrapped with stacked"

	err = stacked.Wrap(stringError("error"))

	err = stringError("error")
	err = stacked.Wrap(err)

	_, err = 0, stringError("error") // want "value converted to error type stringError is not wrapped with stacked"

	_, err = 0, stacked.Wrap(stringError("error"))

	_, err = 0, stringError("error")
	err = stacked.Wrap(err)
}

func stringLiteralDeclarationSelf() {
	{
		var err = stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(stringError("error"))
		_ = err
	}

	{
		var err error = stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
		_ = err
	}

	{
		var err error = stringError("error")
		err = stacked.Wrap(err)
	}

	{
		var err = error(stringError("error")) // want "value converted to error type stringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(stringError("error"))
		_ = err
	}

	{
		var _, err = 0, error(stringError("error")) // want "value converted to error type stringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralShortDeclarationSelf() {
	{
		err := stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(stringError("error"))
		_ = err
	}

	{
		err := error(stringError("error")) // want "value converted to error type stringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(stringError("error"))
		_ = err
	}

	{
		_, err := 0, error(stringError("error")) // want "value converted to error type stringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralReturnSingleSelf() error {
	return stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
	return stacked.Wrap(stringError("error"))
}

func stringLiteralReturnMultipleSelf() (int, error) {
	return 0, stringError("error") // want "value converted to error type stringError is not wrapped with stacked"
	return 0, stacked.Wrap(stringError("error"))
}

func stringLiteralArgumentSelf() {
	functionWithIntErrorArgument(0, stringError("error")) // want "value converted to error type stringError is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(stringError("error")))
}

func stringLiteralCompositeLiteralSelf() {
	_ = structWithErrorField{
		err: stringError("error"), // want "value converted to error type stringError is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(stringError("error")),
	}

	_ = []error{stringError("error")} // want "value converted to error type stringError is not wrapped with stacked"
	_ = []error{stacked.Wrap(stringError("error"))}

	_ = map[string]error{"": stringError("error")} // want "value converted to error type stringError is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(stringError("error"))}
}
