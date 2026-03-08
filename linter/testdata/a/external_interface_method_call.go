package a

import (
	"io/fs"

	"github.com/tbeati/stacked"
)

func externalInterfaceMethodCallAssignmentExternal() {
	var err error
	_ = err
	var file fs.File

	err = file.Close() // want "error returned by file.Close is not wrapped with stacked"
	err = stacked.Wrap(file.Close())

	_, err = 0, file.Close() // want "error returned by file.Close is not wrapped with stacked"
	_, err = 0, stacked.Wrap(file.Close())

	_, err = file.Stat() // want "error returned by file.Stat is not wrapped with stacked"
	_, err = stacked.Wrap2(file.Stat())
}

func externalInterfaceMethodCallDeclarationExternal() {
	var file fs.File

	{
		var err = file.Close() // want "error returned by file.Close is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(file.Close())
		_ = err
	}

	{
		var _, err = 0, file.Close() // want "error returned by file.Close is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(file.Close())
		_ = err
	}

	{
		var _, err = file.Stat() // want "error returned by file.Stat is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = stacked.Wrap2(file.Stat())
		_ = err
	}
}

func externalInterfaceMethodCallShortDeclarationExternal() {
	var file fs.File

	{
		err := file.Close() // want "error returned by file.Close is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(file.Close())
		_ = err
	}

	{
		_, err := 0, file.Close() // want "error returned by file.Close is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(file.Close())
		_ = err
	}

	{
		_, err := file.Stat() // want "error returned by file.Stat is not wrapped with stacked"
		_ = err
	}
	{
		_, err := stacked.Wrap2(file.Stat())
		_ = err
	}
}

func externalInterfaceMethodCallReturn1External() error {
	var file fs.File

	return file.Close() // want "error returned by file.Close is not wrapped with stacked"
	return stacked.Wrap(file.Close())
}

func externalInterfaceMethodCallReturn2External() (fs.FileInfo, error) {
	var file fs.File

	return nil, file.Close() // want "error returned by file.Close is not wrapped with stacked"
	return nil, stacked.Wrap(file.Close())

	return file.Stat() // want "error returned by file.Stat is not wrapped with stacked"
	return stacked.Wrap2(file.Stat())
}

func externalInterfaceMethodCallArgumentExternal() {
	var file fs.File

	functionWithErrorArgument(file.Close()) // want "error returned by file.Close is not wrapped with stacked"
	functionWithErrorArgument(stacked.Wrap(file.Close()))

	functionWithFileErrorArgument(nil, file.Close()) // want "error returned by file.Close is not wrapped with stacked"
	functionWithFileErrorArgument(nil, stacked.Wrap(file.Close()))

	functionWithFileInfoErrorArgument(file.Stat()) // want "error returned by file.Stat is not wrapped with stacked"
	functionWithFileInfoErrorArgument(stacked.Wrap2(file.Stat()))
}

func externalInterfaceMethodCallCompositeLiteralExternal() {
	var file fs.File

	_ = structWithErrorField{
		err: file.Close(), // want "error returned by file.Close is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(file.Close()),
	}

	_ = []error{file.Close()} // want "error returned by file.Close is not wrapped with stacked"
	_ = []error{stacked.Wrap(file.Close())}

	_ = map[string]error{"": file.Close()} // want "error returned by file.Close is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(file.Close())}
}

func externalInterfaceMethodCallChannelSendExternal() {
	var errChan chan error
	var file fs.File

	errChan <- file.Close() // want "error returned by file.Close is not wrapped with stacked"
	errChan <- stacked.Wrap(file.Close())
}

func externalInterfaceMethodCallBlankAssignmentExternal() {
	var file fs.File

	_ = file.Close()
	_, _ = 0, file.Close()
	_, _ = file.Stat()
}
