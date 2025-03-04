package a

import (
	"github.com/tbeati/stacked"

	"testdata/generated"
)

func methodCallAssignmentGenerated() error {
	var err error
	s := generated.StructWithMethods{}

	err = s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = stacked.Wrap(s.SingleReturn())
	if err != nil {
		return err
	}

	err = s.SingleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(s.SingleReturn())
	if err != nil {
		return err
	}

	_, err = 0, s.SingleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = s.MultipleReturn() // want "error returned by s.MultipleReturn is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = s.MultipleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func methodCallDeclarationGenerated() error {
	s := generated.StructWithMethods{}

	{
		var err = s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(s.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		var err = s.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(s.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, s.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = s.MultipleReturn() // want "error returned by s.MultipleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = s.MultipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallShortDeclarationGenerated() error {
	s := generated.StructWithMethods{}

	{
		err := s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(s.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		err := s.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(s.SingleReturn())
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, s.SingleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := s.MultipleReturn() // want "error returned by s.MultipleReturn is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := s.MultipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallReturnSingleGenerated() error {
	s := generated.StructWithMethods{}

	return s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
	return stacked.Wrap(s.SingleReturn())
}

func methodCallReturnMultipleGenerated() (int, error) {
	s := generated.StructWithMethods{}

	return 0, s.SingleReturn() // want "error returned by s.SingleReturn is not wrapped with stacked"
	return 0, stacked.Wrap(s.SingleReturn())
	return s.MultipleReturn() // want "error returned by s.MultipleReturn is not wrapped with stacked"
}

func methodCallArgumentGenerated() {
	s := generated.StructWithMethods{}

	functionWithIntErrorArgument(0, s.SingleReturn()) // want "error returned by s.SingleReturn is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(s.SingleReturn()))
	functionWithIntErrorArgument(s.MultipleReturn()) // want "error returned by s.MultipleReturn is not wrapped with stacked"
}

func methodCallCompositeLiteralGenerated() {
	s := generated.StructWithMethods{}

	_ = structWithErrorField{
		err: s.SingleReturn(), // want "error returned by s.SingleReturn is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(s.SingleReturn()),
	}

	_ = []error{s.SingleReturn()} // want "error returned by s.SingleReturn is not wrapped with stacked"
	_ = []error{stacked.Wrap(s.SingleReturn())}

	_ = map[string]error{"": s.SingleReturn()} // want "error returned by s.SingleReturn is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(s.SingleReturn())}
}

func methodCallIgnoredGenerated() {
	var err error
	s := generated.IgnoredStruct{}
	err = s.IgnoredMethod(err)
	err = (&s).IgnoredMethod(err)
}
