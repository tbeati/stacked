package a

import (
	"github.com/beati/stacked"

	"testdata/b"
)

func structLiteralAssignmentInternal() {
	var err error

	err = b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"

	err = stacked.Wrap(b.StructError{Message: "error"})

	err = b.StructError{Message: "error"}
	err = stacked.Wrap(err)

	_, err = 0, b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"

	_, err = 0, stacked.Wrap(b.StructError{Message: "error"})

	_, err = 0, b.StructError{Message: "error"}
	err = stacked.Wrap(err)
}

func structLiteralDeclarationInternal() {
	{
		var err = b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err = stacked.Wrap(b.StructError{Message: "error"})
		_ = err
	}

	{
		var err error = b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err error = b.StructError{Message: "error"}
		err = stacked.Wrap(err)
	}

	{
		var err = error(b.StructError{Message: "error"}) // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var err = error(b.StructError{Message: "error"})
		err = stacked.Wrap(err)
	}

	{
		var _, err = 0, b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, stacked.Wrap(b.StructError{Message: "error"})
		_ = err
	}

	{
		var _, err = 0, error(b.StructError{Message: "error"}) // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		var _, err = 0, error(b.StructError{Message: "error"})
		err = stacked.Wrap(err)
	}
}

func structLiteralShortDeclarationInternal() {
	{
		err := b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		err := stacked.Wrap(b.StructError{Message: "error"})
		_ = err
	}

	{
		err := error(b.StructError{Message: "error"}) // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		err := error(b.StructError{Message: "error"})
		err = stacked.Wrap(err)
	}

	{
		_, err := 0, b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, stacked.Wrap(b.StructError{Message: "error"})
		_ = err
	}

	{
		_, err := 0, error(b.StructError{Message: "error"}) // want "b.StructError literal is not wrapped with stacked"
		_ = err
	}

	{
		_, err := 0, error(b.StructError{Message: "error"})
		err = stacked.Wrap(err)
	}
}

func structLiteralReturnSingleInternal() error {
	return b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
	return stacked.Wrap(b.StructError{Message: "error"})
}

func structLiteralReturnMultipleInternal() (int, error) {
	return 0, b.StructError{Message: "error"} // want "b.StructError literal is not wrapped with stacked"
	return 0, stacked.Wrap(b.StructError{Message: "error"})
}

func structLiteralArgumentInternal() {
	errArgument(0, b.StructError{Message: "error"}) // want "b.StructError literal is not wrapped with stacked"
	errArgument(0, stacked.Wrap(b.StructError{Message: "error"}))
}

func structLiteralCompositeLiteralInternal() {
	_ = errStruct{
		err: b.StructError{Message: "error"}, // want "b.StructError literal is not wrapped with stacked"
	}
	_ = errStruct{
		err: stacked.Wrap(b.StructError{Message: "error"}),
	}

	_ = []error{b.StructError{Message: "error"}} // want "b.StructError literal is not wrapped with stacked"
	_ = []error{stacked.Wrap(b.StructError{Message: "error"})}

	_ = map[string]error{"": b.StructError{Message: "error"}} // want "b.StructError literal is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(b.StructError{Message: "error"})}
}
