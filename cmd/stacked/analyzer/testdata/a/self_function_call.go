package a

import (
	"github.com/beati/stacked"
)

func functionCallAssignmentSelf() error {
	var err error

	err = singleReturn()
	if err != nil {
		return err
	}

	err = stacked.Wrap(singleReturn())
	if err != nil {
		return err
	}

	err = singleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, singleReturn()
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(singleReturn())
	if err != nil {
		return err
	}

	_, err = 0, singleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = multipleReturn()
	if err != nil {
		return err
	}

	_, err = multipleReturn()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func functionCallDeclarationSelf() error {
	{
		var err = singleReturn()
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(singleReturn())
		if err != nil {
			return err
		}
	}

	{
		var err = singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, singleReturn()
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(singleReturn())
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = multipleReturn()
		if err != nil {
			return err
		}
	}

	{
		var _, err = multipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallShortDeclarationSelf() error {
	{
		err := singleReturn()
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(singleReturn())
		if err != nil {
			return err
		}
	}

	{
		err := singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, singleReturn()
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(singleReturn())
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, singleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := multipleReturn()
		if err != nil {
			return err
		}
	}

	{
		_, err := multipleReturn()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallReturnSingleSelf() error {
	return singleReturn()
	return stacked.Wrap(singleReturn())
}

func functionCallReturnMultipleSelf() (int, error) {
	return 0, singleReturn()
	return 0, stacked.Wrap(singleReturn())
	return multipleReturn()
}

func functionCallArgumentSelf() {
	functionWithIntErrorArgument(0, singleReturn())
	functionWithIntErrorArgument(0, stacked.Wrap(singleReturn()))
	functionWithIntErrorArgument(multipleReturn())
}

func functionCallCompositeLiteralSelf() {
	_ = structWithErrorField{
		err: singleReturn(),
	}
	_ = structWithErrorField{
		err: stacked.Wrap(singleReturn()),
	}

	_ = []error{singleReturn()}
	_ = []error{stacked.Wrap(singleReturn())}

	_ = map[string]error{"": singleReturn()}
	_ = map[string]error{"": stacked.Wrap(singleReturn())}
}
