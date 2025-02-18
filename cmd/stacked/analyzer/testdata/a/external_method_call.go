package a

import (
	"os"

	"github.com/beati/stacked"
)

func methodCallAssignmentExternal() error {
	var err error
	var file *os.File

	err = file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = stacked.Wrap(file.Chdir())
	if err != nil {
		return err
	}

	err = file.Chdir()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(file.Chdir())
	if err != nil {
		return err
	}

	_, err = 0, file.Chdir()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = file.Read(nil) // want "error returned by file.Read is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = file.Read(nil)
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func methodCallDeclarationExternal() error {
	var file *os.File

	{
		var err = file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(file.Chdir())
		if err != nil {
			return err
		}
	}

	{
		var err = file.Chdir()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(file.Chdir())
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, file.Chdir()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = file.Read(nil) // want "error returned by file.Read is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = file.Read(nil)
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallShortDeclarationExternal() error {
	var file *os.File

	{
		err := file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(file.Chdir())
		if err != nil {
			return err
		}
	}

	{
		err := file.Chdir()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(file.Chdir())
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, file.Chdir()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := file.Read(nil) // want "error returned by file.Read is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := file.Read(nil)
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func methodCallReturnSingleExternal() error {
	var file *os.File

	return file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
	return stacked.Wrap(file.Chdir())
}

func methodCallReturnMultipleExternal() (int, error) {
	var file *os.File

	return 0, file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
	return 0, stacked.Wrap(file.Chdir())
	return file.Read(nil) // want "error returned by file.Read is not wrapped with stacked"
}

func methodCallArgumentExternal() {
	var file *os.File

	functionWithIntErrorArgument(0, file.Chdir()) // want "error returned by file.Chdir is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(file.Chdir()))
	functionWithIntErrorArgument(file.Read(nil)) // want "error returned by file.Read is not wrapped with stacked"
}

func methodCallCompositeLiteralExternal() {
	var file *os.File

	_ = structWithErrorField{
		err: file.Chdir(), // want "error returned by file.Chdir is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(file.Chdir()),
	}

	_ = []error{file.Chdir()} // want "error returned by file.Chdir is not wrapped with stacked"
	_ = []error{stacked.Wrap(file.Chdir())}

	_ = map[string]error{"": file.Chdir()} // want "error returned by file.Chdir is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(file.Chdir())}
}
