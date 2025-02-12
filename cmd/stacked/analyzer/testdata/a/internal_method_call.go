package a

import (
	"github.com/beati/stacked"

	"testdata/b"
)

func methodCallAssignmentInternal() error {
	var err error
	s := b.S{}

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
	s := b.S{}

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
	s := b.S{}

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
	s := b.S{}

	return s.SingleReturn()
	return stacked.Wrap(s.SingleReturn())
}

func methodCallReturnMultipleInternal() (int, error) {
	s := b.S{}

	return 0, s.SingleReturn()
	return 0, stacked.Wrap(s.SingleReturn())
	return s.MultipleReturn()
}

func methodCallArgumentInternal() {
	s := b.S{}

	errArgument(0, s.SingleReturn())
	errArgument(0, stacked.Wrap(s.SingleReturn()))
	errArgument(s.MultipleReturn())
}

func methodCallCompositeLiteralInternal() {
	s := b.S{}

	_ = errStruct{
		err: s.SingleReturn(),
	}
	_ = errStruct{
		err: stacked.Wrap(s.SingleReturn()),
	}

	_ = []error{s.SingleReturn()}
	_ = []error{stacked.Wrap(s.SingleReturn())}

	_ = map[string]error{"": s.SingleReturn()}
	_ = map[string]error{"": stacked.Wrap(s.SingleReturn())}
}
