# stacked

[![Go Reference](https://pkg.go.dev/badge/github.com/tbeati/stacked.svg)](https://pkg.go.dev/github.com/tbeati/stacked)
[![License](https://img.shields.io/github/license/tbeati/stacked)](LICENSE)

`stacked` is a Go library that attaches stack traces to your errors right at their source. While standard Go error handling often leaves you guessing where an issue actually originated as it bubbles up through intermediate functions, `stacked` captures the context the moment the error is produced.

When paired with its dedicated linter, `stacked` allows you to enforce a strict "wrap at the source" policy across your entire codebase, offering a seamless debugging experience:

* **Zero Guesswork:** By capturing the exact function, file, and line number where the error occurred, it cuts down debugging time drastically.
* **Foolproof Coverage:** The linter acts as a safety net, guaranteeing that no error is left unwrapped before your code even compiles.
* **Frictionless Integration:** Wrapping is idempotent (the first wrap wins) and fully compatible with the standard library (`errors.Is`, `errors.AsType`, and `errors.Unwrap`), meaning your existing error-handling logic remains completely intact.

Beyond standard errors, `stacked` provides a `Recover` utility to catch panics and convert them into `stacked` errors. This allows you to log the panic's stack trace using your own custom format.

## Install

```sh
go get github.com/tbeati/stacked
```

## Usage

### Wrap and log

Wrap the call that produces the error, then log it with `slog`, attaching the
error and its stack trace as structured fields:

```go
// Wrap at the source: the trace is captured here, pinning os.Chdir.
err := stacked.Wrap(os.Chdir("/no/such/directory"))
if err != nil {
	slog.Error("failed to change directory",
		slog.Any("error", err),
		slog.Any("stack", stacked.StackTrace(err)),
	)
	return
}
```

`stacked.StackTrace(err)` returns `[]stacked.StackFrame`, and each frame has
`Function`, `File`, and `Line` fields with JSON tags — so with
`slog.NewJSONHandler` the trace serializes as a clean array of frames ready
for your log pipeline.

### Higher arity

`Wrap2`–`Wrap5` cover calls that return extra values alongside the error —
the result passes straight through:

```go
data, err := stacked.Wrap2(os.ReadFile("/no/such/file"))
```

### Iterators

For range-over-func iterators, wrap the sequence with `WrapSeq` (or
`WrapSeq2` for `iter.Seq2[T, error]`) — the trace is captured at the yield
site:

```go
for err := range stacked.WrapSeq(produceErrors()) {
	slog.Error("step failed",
		slog.Any("error", err),
		slog.Any("stack", stacked.StackTrace(err)),
	)
}
```

When you drive a sequence with `iter.Pull`, wrap each pull with `WrapPull`
(or `WrapPull2`) so the trace points at the `next()` call site instead:

```go
next, stop := iter.Pull(produceErrors())
defer stop()
for {
	err, ok := stacked.WrapPull(next())
	if !ok {
		break
	}
	slog.Error("step failed", slog.Any("stack", stacked.StackTrace(err)))
}
```

### Reformat panic stack traces

`Recover` runs a function and turns any panic (or `runtime.Goexit`) into a
`stacked` error, with the trace pinned at the panic site:

```go
stacked.Recover(
    func() {
		// code to run
    },
	func(err error) {
        slog.Error("panic occured",
            slog.Any("panic", true),
            slog.Any("error", err),
            slog.Any("stack", stacked.StackTrace(err)),
        )
    },
	true,
)
```

### Ignore sentinel errors

Some errors are expected control flow, not failures, and shouldn't carry a
trace. Register them so every `Wrap` call returns them untouched
(`io.EOF` is registered by default):

```go
// By value: future Wrap calls return this error — or any error that
// wraps it (matched with errors.Is) — unchanged.
stacked.Ignore(sql.ErrNoRows)

// By predicate: ignore a whole class of errors by type.
stacked.IgnoreFunc(func(err error) bool {
	_, ok := errors.AsType[*os.PathError](err)
	return ok
})

// Now wrapping these is a no-op — no stack trace attached.
err := stacked.Wrap(row.Scan(&v)) // sql.ErrNoRows passes straight through
```
