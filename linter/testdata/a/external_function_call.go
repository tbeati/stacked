package a

import (
	"os"

	"github.com/tbeati/stacked"
)

func functionCallAssignmentExternal() error {
	var err error

	err = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = stacked.Wrap(os.Chdir("/"))
	if err != nil {
		return err
	}

	err = os.Chdir("/")
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = 0, stacked.Wrap(os.Chdir("/"))
	if err != nil {
		return err
	}

	_, err = 0, os.Chdir("/")
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	_, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	if err != nil {
		return err
	}

	_, err = os.Hostname()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	return nil
}

func functionCallDeclarationExternal() error {
	{
		var err = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var err = stacked.Wrap(os.Chdir("/"))
		if err != nil {
			return err
		}
	}

	{
		var err = os.Chdir("/")
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, stacked.Wrap(os.Chdir("/"))
		if err != nil {
			return err
		}
	}

	{
		var _, err = 0, os.Chdir("/")
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		var _, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		var _, err = os.Hostname()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallShortDeclarationExternal() error {
	{
		err := os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		err := stacked.Wrap(os.Chdir("/"))
		if err != nil {
			return err
		}
	}

	{
		err := os.Chdir("/")
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, stacked.Wrap(os.Chdir("/"))
		if err != nil {
			return err
		}
	}

	{
		_, err := 0, os.Chdir("/")
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	{
		_, err := os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
		if err != nil {
			return err
		}
	}

	{
		_, err := os.Hostname()
		err = stacked.Wrap(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func functionCallReturnSingleExternal() error {
	return os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	return stacked.Wrap(os.Chdir("/"))
}

func functionCallReturnMultipleExternal() (string, error) {
	return "", os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	return "", stacked.Wrap(os.Chdir("/"))
	return os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	return stacked.Wrap2(os.Hostname())
}

func functionCallArgumentExternal() {
	functionWithStringErrorArgument("", os.Chdir("/")) // want "error returned by os.Chdir is not wrapped with stacked"
	functionWithStringErrorArgument("", stacked.Wrap(os.Chdir("/")))
	functionWithStringErrorArgument(os.Hostname()) // want "error returned by os.Hostname is not wrapped with stacked"
}

func functionCallCompositeLiteralExternal() {
	_ = structWithErrorField{
		err: os.Chdir("/"), // want "error returned by os.Chdir is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(os.Chdir("/")),
	}

	_ = []error{os.Chdir("/")} // want "error returned by os.Chdir is not wrapped with stacked"
	_ = []error{stacked.Wrap(os.Chdir("/"))}

	_ = map[string]error{"": os.Chdir("/")} // want "error returned by os.Chdir is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(os.Chdir("/"))}
}
