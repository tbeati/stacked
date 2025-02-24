package a

import (
	"os"

	"github.com/tbeati/stacked"
)

/*
	{
		var x, y = os.Open("")
		_, _ = x, y

		var (
			a, b     = os.Open("")
			c, d     = 2, ""
			e, f int = 2, 4
		)

		_, _ = a, b
		_, _ = c, d
		_, _ = e, f
	}
*/

var foo = structWithErrorField{
	err: os.Chdir("/"),
}

func callExternalPackage() error {
	var name string
	var err2 error

	var err3 = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
	if err3 != nil {
		return err3
	}

	name, err := os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	err = stacked.Wrap(err2)
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	err2 = stacked.Wrap(err)
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "error returned by os.Hostname is not wrapped with stacked"
	err2 = stacked.Wrap(err2)
	if err != nil {
		return err
	}

	name, err = os.Hostname()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	if err := os.Chmod(name, 0777); err != nil { // want "error returned by os.Chmod is not wrapped with stacked"
		return err
	}

	(err) = (os.Chmod(name, 0777)) // want "error returned by os.Chmod is not wrapped with stacked"

	err = os.Chmod(name, 0777) // want "error returned by os.Chmod is not wrapped with stacked"
	name = "test"
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	err = stacked.Wrap(os.Chmod("test", 0777))
	if err != nil {
		return err
	}

	err = os.Chmod("test", 0777)
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	err = os.Chmod("test", 0777) // want "error returned by os.Chmod is not wrapped with stacked"
	if err != nil {
		return stacked.Wrap(err)
	}

	f, err := os.Open(name) // want "error returned by os.Open is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = f.Close() // want "error returned by f.Close is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = stacked.Wrap(f.Close())
	if err != nil {
		return err
	}

	err = f.Close()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	es := structWithErrorField{
		err: f.Close(), // want "error returned by f.Close is not wrapped with stacked"
	}

	es = structWithErrorField{
		f.Close(), // want "error returned by f.Close is not wrapped with stacked"
	}

	es.err = f.Close() // want "error returned by f.Close is not wrapped with stacked"
	if err != nil {
		return err
	}

	errSlice := []error{
		f.Close(), // want "error returned by f.Close is not wrapped with stacked"
	}

	errSlice = append(errSlice, f.Close()) // want "error returned by f.Close is not wrapped with stacked"

	errSlice[0] = f.Close()   // want "error returned by f.Close is not wrapped with stacked"
	errSlice[1+1] = f.Close() // want "error returned by f.Close is not wrapped with stacked"

	errSlice[1+1] = f.Close()
	errSlice[1+1] = stacked.Wrap(errSlice[1+1])

	errMap := map[int]error{
		0: f.Close(), // want "error returned by f.Close is not wrapped with stacked"
	}

	errMap[1] = f.Close() // want "error returned by f.Close is not wrapped with stacked"

	var errPointer *error
	*errPointer = f.Close() // want "error returned by f.Close is not wrapped with stacked"

	return f.Close() // want "error returned by f.Close is not wrapped with stacked"
}
