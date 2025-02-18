package a

import (
	"github.com/tbeati/stacked"

	"testdata/b"
)

func methodCallAssignmentInternal() error {
	var err error
	s := b.StructWithMethods{}

	err = s.SingleReturn()
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

	_, err = 0, s.SingleReturn()
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

	_, err = s.MultipleReturn()
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

func methodCallDeclarationInternal() error {
	s := b.StructWithMethods{}

	{
		var err = s.SingleReturn()
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
		var _, err = 0, s.SingleReturn()
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
		var _, err = s.MultipleReturn()
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

func methodCallShortDeclarationInternal() error {
	s := b.StructWithMethods{}

	{
		err := s.SingleReturn()
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
		_, err := 0, s.SingleReturn()
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
		_, err := s.MultipleReturn()
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

func methodCallReturnSingleInternal() error {
	s := b.StructWithMethods{}

	return s.SingleReturn()
	return stacked.Wrap(s.SingleReturn())
}

func methodCallReturnMultipleInternal() (int, error) {
	s := b.StructWithMethods{}

	return 0, s.SingleReturn()
	return 0, stacked.Wrap(s.SingleReturn())
	return s.MultipleReturn()
}

func methodCallArgumentInternal() {
	s := b.StructWithMethods{}

	functionWithIntErrorArgument(0, s.SingleReturn())
	functionWithIntErrorArgument(0, stacked.Wrap(s.SingleReturn()))
	functionWithIntErrorArgument(s.MultipleReturn())
}

func methodCallCompositeLiteralInternal() {
	s := b.StructWithMethods{}

	_ = structWithErrorField{
		err: s.SingleReturn(),
	}
	_ = structWithErrorField{
		err: stacked.Wrap(s.SingleReturn()),
	}

	_ = []error{s.SingleReturn()}
	_ = []error{stacked.Wrap(s.SingleReturn())}

	_ = map[string]error{"": s.SingleReturn()}
	_ = map[string]error{"": stacked.Wrap(s.SingleReturn())}
}
