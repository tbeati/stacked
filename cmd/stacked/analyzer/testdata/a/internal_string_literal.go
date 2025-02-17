package a

import (
	"github.com/beati/stacked"

	"testdata/b"
)

func stringLiteralAssignmentInternal() {
	var err error

	err = b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"

	err = stacked.Wrap(b.StringError("error"))

	err = b.StringError("error")
	err = stacked.Wrap(err)

	_, err = 0, b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"

	_, err = 0, stacked.Wrap(b.StringError("error"))

	_, err = 0, b.StringError("error")
	err = stacked.Wrap(err)
}

func stringLiteralDeclarationInternal() {
	{
		var err = b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(b.StringError("error"))
		_ = err
	}

	{
		var err error = b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var err error = b.StringError("error")
		err = stacked.Wrap(err)
	}

	{
		var err = error(b.StringError("error")) // want "value converted to error type b.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(b.StringError("error"))
		_ = err
	}

	{
		var _, err = 0, error(b.StringError("error")) // want "value converted to error type b.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralShortDeclarationInternal() {
	{
		err := b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(b.StringError("error"))
		_ = err
	}

	{
		err := error(b.StringError("error")) // want "value converted to error type b.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(b.StringError("error"))
		_ = err
	}

	{
		_, err := 0, error(b.StringError("error")) // want "value converted to error type b.StringError is not wrapped with stacked" "value converted to error type error is not wrapped with stacked"
		_ = err
	}
}

func stringLiteralReturnSingleInternal() error {
	return b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
	return stacked.Wrap(b.StringError("error"))
}

func stringLiteralReturnMultipleInternal() (int, error) {
	return 0, b.StringError("error") // want "value converted to error type b.StringError is not wrapped with stacked"
	return 0, stacked.Wrap(b.StringError("error"))
}

func stringLiteralArgumentInternal() {
	errArgument(0, b.StringError("error")) // want "value converted to error type b.StringError is not wrapped with stacked"
	errArgument(0, stacked.Wrap(b.StringError("error")))
}

func stringLiteralCompositeLiteralInternal() {
	_ = errStruct{
		err: b.StringError("error"), // want "value converted to error type b.StringError is not wrapped with stacked"
	}
	_ = errStruct{
		err: stacked.Wrap(b.StringError("error")),
	}

	_ = []error{b.StringError("error")} // want "value converted to error type b.StringError is not wrapped with stacked"
	_ = []error{stacked.Wrap(b.StringError("error"))}

	_ = map[string]error{"": b.StringError("error")} // want "value converted to error type b.StringError is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(b.StringError("error"))}
}
