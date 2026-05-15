# stacked

[![Go Reference](https://pkg.go.dev/badge/github.com/tbeati/stacked.svg)](https://pkg.go.dev/github.com/tbeati/stacked)
[![License](https://img.shields.io/github/license/tbeati/stacked)](LICENSE)

`stacked` is a Go library that attaches stack traces to your errors right at their source. While standard Go error handling often leaves you guessing where an issue actually originated as it bubbles up through intermediate functions, `stacked` captures the context the moment the error is produced.

When paired with its dedicated linter, `stacked` allows you to enforce a strict "wrap at the source" policy across your entire codebase, offering a seamless debugging experience:

* **Zero Guesswork:** By capturing the exact function, file, and line number where the error occurred, it cuts down debugging time drastically.
* **Foolproof Coverage:** The linter acts as a safety net, guaranteeing that no error is left unwrapped before your code even compiles.
* **Frictionless Integration:** Wrapping is idempotent (the first wrap wins) and fully compatible with the standard library (`errors.Is`, `errors.As`, and `errors.Unwrap`), meaning your existing error-handling logic remains completely intact.

Beyond standard errors, `stacked` provides a `Recover` utility to catch panics and convert them into `stacked` errors. This allows you to log the panic's stack trace using your own custom format.
