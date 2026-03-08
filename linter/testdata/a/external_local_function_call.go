package a

import (
	"os"

	"github.com/tbeati/stacked"
)

func localFunctionCallAssignmentExternal() {
	localFunc := os.Chdir
	localFunc2 := os.Hostname
	localFunc3 := os.Pipe

	var err error
	_ = err

	err = localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	err = stacked.Wrap(localFunc("/"))

	_, err = 0, localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	_, err = 0, stacked.Wrap(localFunc("/"))

	_, err = localFunc2() // want "error returned by localFunc2 is not wrapped with stacked"
	_, err = stacked.Wrap2(localFunc2())

	_, _, err = localFunc3() // want "error returned by localFunc3 is not wrapped with stacked"
	_, _, err = stacked.Wrap3(localFunc3())
}

func localFunctionCallDeclarationExternal() {
	localFunc := os.Chdir
	localFunc2 := os.Hostname
	localFunc3 := os.Pipe

	{
		var err = localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(localFunc("/"))
		_ = err
	}

	{
		var _, err = 0, localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(localFunc("/"))
		_ = err
	}

	{
		var _, err = localFunc2() // want "error returned by localFunc2 is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = stacked.Wrap2(localFunc2())
		_ = err
	}

	{
		var _, _, err = localFunc3() // want "error returned by localFunc3 is not wrapped with stacked"
		_ = err
	}
	{
		var _, _, err = stacked.Wrap3(localFunc3())
		_ = err
	}
}

func localFunctionCallShortDeclarationExternal() {
	localFunc := os.Chdir
	localFunc2 := os.Hostname
	localFunc3 := os.Pipe

	{
		err := localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(localFunc("/"))
		_ = err
	}

	{
		_, err := 0, localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(localFunc("/"))
		_ = err
	}

	{
		_, err := localFunc2() // want "error returned by localFunc2 is not wrapped with stacked"
		_ = err
	}
	{
		_, err := stacked.Wrap2(localFunc2())
		_ = err
	}

	{
		_, _, err := localFunc3() // want "error returned by localFunc3 is not wrapped with stacked"
		_ = err
	}
	{
		_, _, err := stacked.Wrap3(localFunc3())
		_ = err
	}
}

func localFunctionCallReturn1External() error {
	localFunc := os.Chdir

	return localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	return stacked.Wrap(localFunc("/"))
}

func localFunctionCallReturn2External() (string, error) {
	localFunc := os.Chdir
	localFunc2 := os.Hostname

	return "", localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	return "", stacked.Wrap(localFunc("/"))

	return localFunc2() // want "error returned by localFunc2 is not wrapped with stacked"
	return stacked.Wrap2(localFunc2())
}

func localFunctionCallReturn3External() (*os.File, *os.File, error) {
	localFunc := os.Chdir
	localFunc3 := os.Pipe

	return nil, nil, localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	return nil, nil, stacked.Wrap(localFunc("/"))

	return localFunc3() // want "error returned by localFunc3 is not wrapped with stacked"
	return stacked.Wrap3(localFunc3())
}

func localFunctionCallArgumentExternal() {
	localFunc := os.Chdir
	localFunc2 := os.Hostname
	localFunc3 := os.Pipe

	functionWithErrorArgument(localFunc("/")) // want "error returned by localFunc is not wrapped with stacked"
	functionWithErrorArgument(stacked.Wrap(localFunc("/")))

	functionWithStringErrorArgument("", localFunc("/")) // want "error returned by localFunc is not wrapped with stacked"
	functionWithStringErrorArgument("", stacked.Wrap(localFunc("/")))

	functionWithStringErrorArgument(localFunc2()) // want "error returned by localFunc2 is not wrapped with stacked"
	functionWithStringErrorArgument(stacked.Wrap2(localFunc2()))

	functionWithFileFileErrorArgument(nil, nil, localFunc("/")) // want "error returned by localFunc is not wrapped with stacked"
	functionWithFileFileErrorArgument(nil, nil, stacked.Wrap(localFunc("/")))

	functionWithFileFileErrorArgument(localFunc3()) // want "error returned by localFunc3 is not wrapped with stacked"
	functionWithFileFileErrorArgument(stacked.Wrap3(localFunc3()))
}

func localFunctionCallCompositeLiteralExternal() {
	localFunc := os.Chdir

	_ = structWithErrorField{
		err: localFunc("/"), // want "error returned by localFunc is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(localFunc("/")),
	}

	_ = []error{localFunc("/")} // want "error returned by localFunc is not wrapped with stacked"
	_ = []error{stacked.Wrap(localFunc("/"))}

	_ = map[string]error{"": localFunc("/")} // want "error returned by localFunc is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(localFunc("/"))}
}

func localFunctionCallChannelSendExternal() {
	localFunc := os.Chdir

	var errChan chan error

	errChan <- localFunc("/") // want "error returned by localFunc is not wrapped with stacked"
	errChan <- stacked.Wrap(localFunc("/"))
}

func localFunctionCallBlankAssignmentExternal() {
	localFunc := os.Chdir
	localFunc2 := os.Hostname
	localFunc3 := os.Pipe

	_ = localFunc("/")
	_, _ = 0, localFunc("/")
	_, _ = localFunc2()
	_, _, _ = localFunc3()
}
