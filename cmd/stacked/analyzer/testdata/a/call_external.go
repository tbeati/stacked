package a

import (
	"os"

	"github.com/beati/stacked"
)

func callExternalPackage() error {
	var name string
	var err2 error

	name, err := os.Hostname() // want "err is not wrapped with stacked"
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "err is not wrapped with stacked"
	err = stacked.Wrap(err2)
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "err is not wrapped with stacked"
	err2 = stacked.Wrap(err)
	if err != nil {
		return err
	}

	name, err = os.Hostname() // want "err is not wrapped with stacked"
	err2 = stacked.Wrap(err2)
	if err != nil {
		return err
	}

	name, err = os.Hostname()
	err = stacked.Wrap(err)
	if err != nil {
		return err
	}

	if err := os.Chmod(name, 0777); err != nil { // want "err is not wrapped with stacked"
		return err
	}

	(err) = (os.Chmod(name, 0777)) // want "err is not wrapped with stacked"

	err = os.Chmod(name, 0777) // want "err is not wrapped with stacked"
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

	err = os.Chmod("test", 0777) // want "err is not wrapped with stacked"
	if err != nil {
		return stacked.Wrap(err)
	}

	f, err := os.Open(name) // want "err is not wrapped with stacked"
	if err != nil {
		return err
	}

	err = f.Close() // want "err is not wrapped with stacked"
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

	es := errStruct{
		err: f.Close(), // want "es.err is not wrapped with stacked"
	}

	es.err = f.Close() // want "es.err is not wrapped with stacked"
	if err != nil {
		return err
	}

	errSlice := []error{
		f.Close(), // want "errSlice value is not wrapped with stacked"
	}

	errSlice = append(errSlice, f.Close()) // want "errSlice value is not wrapped with stacked"

	errSlice[0] = f.Close()   // want "errSlice\\[0\\] is not wrapped with stacked"
	errSlice[1+1] = f.Close() // want "errSlice\\[1\\+1\\] is not wrapped with stacked"

	errSlice[1+1] = f.Close()
	errSlice[1+1] = stacked.Wrap(errSlice[1+1])

	errMap := map[int]error{
		0: f.Close(), // want "errMap value is not wrapped with stacked"
	}

	errMap[1] = f.Close() // want "errMap\\[1\\] is not wrapped with stacked"

	var errPointer *error
	*errPointer = f.Close() // want "\\*errPointer is not wrapped with stacked"

	return f.Close() // want "returned error is not wrapped with stacked"
}

type errStruct struct {
	err error
}
