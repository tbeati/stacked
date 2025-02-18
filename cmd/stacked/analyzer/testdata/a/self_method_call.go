package a

import (
	"github.com/tbeati/stacked"
)

func methodCallAssignmentSelf() error {
	var err error
	s := structWithMethods{}

	err = s.singleReturn()
	if err != nil {
		return err
	}

	err = stacked.Wrap(s.singleReturn())
	if err != nil {
		return err
	}

	err = s.singleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, s.singleReturn()
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(s.singleReturn())
	if err != nil {
		return err
	}

	_, err = 0, s.singleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = s.multipleReturn()
	if err != nil {
		return err
	}

	_, err = s.multipleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func methodCallDeclarationSelf() error {
	s := structWithMethods{}

	{
		var err = s.singleReturn()
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(s.singleReturn())
		if err != nil {
			return err
		}
	}

	{
		var err = s.singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, s.singleReturn()
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(s.singleReturn())
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, s.singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = s.multipleReturn()
		if err != nil {
			return err
		}
	}

	{
		var _, err = s.multipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallShortDeclarationSelf() error {
	s := structWithMethods{}

	{
		err := s.singleReturn()
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(s.singleReturn())
		if err != nil {
			return err
		}
	}

	{
		err := s.singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, s.singleReturn()
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(s.singleReturn())
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, s.singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := s.multipleReturn()
		if err != nil {
			return err
		}
	}

	{
		_, err := s.multipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallReturnSingleSelf() error {
	s := structWithMethods{}

	return s.singleReturn()
	return stacked.Wrap(s.singleReturn())
}

func methodCallReturnMultipleSelf() (int, error) {
	s := structWithMethods{}

	return 0, s.singleReturn()
	return 0, stacked.Wrap(s.singleReturn())
	return s.multipleReturn()
}

func methodCallArgumentSelf() {
	s := structWithMethods{}

	functionWithIntErrorArgument(0, s.singleReturn())
	functionWithIntErrorArgument(0, stacked.Wrap(s.singleReturn()))
	functionWithIntErrorArgument(s.multipleReturn())
}

func methodCallCompositeLiteralSelf() {
	s := structWithMethods{}

	_ = structWithErrorField{
		err: s.singleReturn(),
	}
	_ = structWithErrorField{
		err: stacked.Wrap(s.singleReturn()),
	}

	_ = []error{s.singleReturn()}
	_ = []error{stacked.Wrap(s.singleReturn())}

	_ = map[string]error{"": s.singleReturn()}
	_ = map[string]error{"": stacked.Wrap(s.singleReturn())}
}
