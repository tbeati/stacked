package a

import (
	"os"

	"github.com/tbeati/stacked"
)

func edgeCases() {
	{
		var (
			err  = os.Chdir("/") // want "error returned by os.Chdir is not wrapped with stacked"
			_, _ = 0, ""
		)
		_ = err
	}
	{
		var (
			err  = stacked.Wrap(os.Chdir("/"))
			_, _ = 0, ""
		)
		_ = err
	}

	{
		var (
			_, err = os.Open("") // want "error returned by os.Open is not wrapped with stacked"
			_, _   = 0, ""
		)
		_ = err
	}
	{
		var (
			_, err = stacked.Wrap2(os.Open(""))
			_, _   = 0, ""
		)
		_ = err
	}

	if err := os.Chdir("/"); err != nil { // want "error returned by os.Chdir is not wrapped with stacked"
	}
	if err := stacked.Wrap(os.Chdir("/")); err != nil {
	}

	switch err := os.Chdir("/"); err { // want "error returned by os.Chdir is not wrapped with stacked"
	}
	switch err := os.Chdir("/"); err {
	}
}
