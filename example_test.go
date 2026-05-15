package stacked_test

import (
	"fmt"

	"github.com/tbeati/stacked"
)

func ExampleRecover() {
	stacked.Recover(
		func() {
			panic("something went wrong")
		},
		func(err error) {
			fmt.Println(err)

			stackTrace := stacked.StackTrace(err)
			_ = stackTrace // stack trace at the panic site
		},
		false,
	)

	// Output:
	// something went wrong
}
