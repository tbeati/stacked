package a

import (
	"os"

	"github.com/tbeati/stacked"
)

func functionCallAssignmentExternal() {
	var err error
	_ = err

	err = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	err = stacked.Wrap(os.Chdir("/"))

	_, err = 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	_, err = 0, stacked.Wrap(os.Chdir("/"))

	_, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	_, err = stacked.Wrap2(os.Hostname())

	_, _, err = os.Pipe() // want "error returned by os.Pipe is not wrapped with stacked"
	_, _, err = stacked.Wrap3(os.Pipe())
}

func functionCallDeclarationExternal() {
	{
		var err = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(os.Chdir("/"))
		_ = err
	}

	{
		var _, err = 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(os.Chdir("/"))
		_ = err
	}

	{
		var _, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = stacked.Wrap2(os.Hostname())
		_ = err
	}

	{
		var _, _, err = os.Pipe() // want "error returned by os.Pipe is not wrapped with stacked"
		_ = err
	}
	{
		var _, _, err = stacked.Wrap3(os.Pipe())
		_ = err
	}
}

func functionCallShortDeclarationExternal() {
	{
		err := os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(os.Chdir("/"))
		_ = err
	}

	{
		_, err := 0, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(os.Chdir("/"))
		_ = err
	}

	{
		_, err := os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
		_ = err
	}
	{
		_, err := stacked.Wrap2(os.Hostname())
		_ = err
	}

	{
		_, _, err := os.Pipe() // want "error returned by os.Pipe is not wrapped with stacked"
		_ = err
	}
	{
		_, _, err := stacked.Wrap3(os.Pipe())
		_ = err
	}
}

func functionCallReturn1External() error {
	return os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	return stacked.Wrap(os.Chdir("/"))
}

func functionCallReturn2External() (string, error) {
	return "", os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	return "", stacked.Wrap(os.Chdir("/"))

	return os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	return stacked.Wrap2(os.Hostname())
}

func functionCallReturn3External() (*os.File, *os.File, error) {
	return nil, nil, os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	return nil, nil, stacked.Wrap(os.Chdir("/"))

	return os.Pipe() // want "error returned by os.Pipe is not wrapped with stacked"
	return stacked.Wrap3(os.Pipe())
}

func functionCallArgumentExternal() {
	functionWithErrorArgument(os.Chdir("/")) // want "error returned by os.Chdir is not wrapped with stacked"
	functionWithErrorArgument(stacked.Wrap(os.Chdir("/")))

	functionWithStringErrorArgument("", os.Chdir("/")) // want "error returned by os.Chdir is not wrapped with stacked"
	functionWithStringErrorArgument("", stacked.Wrap(os.Chdir("/")))

	functionWithStringErrorArgument(os.Hostname()) // want "error returned by os.Hostname is not wrapped with stacked"
	functionWithStringErrorArgument(stacked.Wrap2(os.Hostname()))

	functionWithFileFileErrorArgument(nil, nil, os.Chdir("/")) // want "error returned by os.Chdir is not wrapped with stacked"
	functionWithFileFileErrorArgument(nil, nil, stacked.Wrap(os.Chdir("/")))

	functionWithFileFileErrorArgument(os.Pipe()) // want "error returned by os.Pipe is not wrapped with stacked"
	functionWithFileFileErrorArgument(stacked.Wrap3(os.Pipe()))
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

func functionCallChannelSendExternal() {
	var errChan chan error

	errChan <- os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	errChan <- stacked.Wrap(os.Chdir("/"))
}

func functionCallBlankAssignmentExternal() {
	_ = os.Chdir("/")
	_, _ = 0, os.Chdir("/")
	_, _ = os.Hostname()
	_, _, _ = os.Pipe()
}
