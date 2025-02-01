package a

import (
	"os"

	"github.com/beati/stacked"
)

func f() error {
	name, err := os.Hostname() // want "coucou"
	if err != nil {
		return err
	}
	_ = name

	if err := os.Chmod("test", 0777); err != nil { // want "coucou"
		return err
	}

	err = stacked.Wrap(os.Chmod("test", 0777))
	if err != nil {
		return err
	}

	err = os.Chmod("test", 0777)
	if err != nil {
		return stacked.Wrap(err)
	}

	err = g()
	if err != nil {
		return err
	}

	return nil
}

func g() error {
	return nil
}
