package a

import (
	"fmt"

	"github.com/tbeati/stacked"

	"testdata/generated"
)

func functionCallAssignmentGenerated() error {
	var err error

	err = generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = stacked.Wrap(generated.SingleReturn())
	if err != nil {
		return err
	}

	err = generated.SingleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(generated.SingleReturn())
	if err != nil {
		return err
	}

	_, err = 0, generated.SingleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = generated.MultipleReturn() // want "error returned by generated.MultipleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = generated.MultipleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func functionCallDeclarationGenerated() error {
	{
		var err = generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(generated.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		var err = generated.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(generated.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, generated.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = generated.MultipleReturn() // want "error returned by generated.MultipleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = generated.MultipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallShortDeclarationGenerated() error {
	{
		err := generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(generated.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		err := generated.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(generated.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, generated.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := generated.MultipleReturn() // want "error returned by generated.MultipleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := generated.MultipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallReturnSingleGenerated() error {
	return generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
	return stacked.Wrap(generated.SingleReturn())
	return generated.ReturnConcreteType() // want "error returned by generated.ReturnConcreteType is not wrapped with stacked"
	return stacked.Wrap(generated.ReturnConcreteType())
	return generated.ReturnConcreteTypePointer() // want "error returned by generated.ReturnConcreteTypePointer is not wrapped with stacked"
	return stacked.Wrap(generated.ReturnConcreteTypePointer())
}

func functionCallReturnMultipleGenerated() (int, error) {
	return 0, generated.SingleReturn() // want "error returned by generated.SingleReturn is not wrapped with stacked"
	return 0, stacked.Wrap(generated.SingleReturn())
	return generated.MultipleReturn() // want "error returned by generated.MultipleReturn is not wrapped with stacked"
}

func functionCallArgumentGenerated() {
	functionWithIntErrorArgument(0, generated.SingleReturn()) // want "error returned by generated.SingleReturn is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(generated.SingleReturn()))
	functionWithIntErrorArgument(generated.MultipleReturn()) // want "error returned by generated.MultipleReturn is not wrapped with stacked"
}

func functionCallCompositeLiteralGenerated() {
	_ = structWithErrorField{
		err: generated.SingleReturn(), // want "error returned by generated.SingleReturn is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(generated.SingleReturn()),
	}

	_ = []error{generated.SingleReturn()} // want "error returned by generated.SingleReturn is not wrapped with stacked"
	_ = []error{stacked.Wrap(generated.SingleReturn())}

	_ = map[string]error{"": generated.SingleReturn()} // want "error returned by generated.SingleReturn is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(generated.SingleReturn())}
}

func functionCallIgnoredGenerated() {
	var err error
	err = generated.IgnoredFunction(err)
	err = fmt.Errorf("wrapping %w", err)
	err = fmt.Errorf("not wrapping") // want "error returned by fmt.Errorf is not wrapped with stacked"
}

func functionCallConcreteTypeGenerated() {
	var err error
	err = generated.ReturnConcreteType()        // want "error returned by generated.ReturnConcreteType is not wrapped with stacked"
	err = generated.ReturnConcreteTypePointer() // want "error returned by generated.ReturnConcreteTypePointer is not wrapped with stacked"
	_ = err
}
