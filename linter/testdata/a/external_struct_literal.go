package a

import (
	"io/fs"

	"github.com/tbeati/stacked"
)

func structLiteralAssignmentExternal() {
	var err error
	_ = err

	err = &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
	err = stacked.Wrap(&fs.PathError{Path: "error"})

	_, err = 0, &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
	_, err = 0, stacked.Wrap(&fs.PathError{Path: "error"})
}

func structLiteralDeclarationExternal() {
	{
		var err = &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(&fs.PathError{Path: "error"})
		_ = err
	}

	{
		var _, err = 0, &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(&fs.PathError{Path: "error"})
		_ = err
	}
}

func structLiteralShortDeclarationExternal() {
	{
		err := &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(&fs.PathError{Path: "error"})
		_ = err
	}

	{
		_, err := 0, &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(&fs.PathError{Path: "error"})
		_ = err
	}
}

func structLiteralReturnSingleExternal() error {
	return &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
	return stacked.Wrap(&fs.PathError{Path: "error"})
}

func structLiteralReturnMultipleExternal() (int, error) {
	return 0, &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
	return 0, stacked.Wrap(&fs.PathError{Path: "error"})
}

func structLiteralArgumentExternal() {
	functionWithIntErrorArgument(0, &fs.PathError{Path: "error"}) // want "fs.PathError literal is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(&fs.PathError{Path: "error"}))
}

func structLiteralCompositeLiteralExternal() {
	_ = structWithErrorField{
		err: &fs.PathError{Path: "error"}, // want "fs.PathError literal is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(&fs.PathError{Path: "error"}),
	}

	_ = []error{&fs.PathError{Path: "error"}} // want "fs.PathError literal is not wrapped with stacked"
	_ = []error{stacked.Wrap(&fs.PathError{Path: "error"})}

	_ = map[string]error{"": &fs.PathError{Path: "error"}} // want "fs.PathError literal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(&fs.PathError{Path: "error"})}
}

func structLiteralChannelSendExternal() {
	var errChan chan error

	errChan <- &fs.PathError{Path: "error"} // want "fs.PathError literal is not wrapped with stacked"
	errChan <- stacked.Wrap(&fs.PathError{Path: "error"})
}
