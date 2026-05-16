# stacked

[![Go Reference](https://pkg.go.dev/badge/github.com/tbeati/stacked.svg)](https://pkg.go.dev/github.com/tbeati/stacked)
[![License](https://img.shields.io/github/license/tbeati/stacked)](LICENSE)

`stacked` is a Go library that attaches stack traces to your errors right at their source, with a linter that enforces wrapping at every error site. While standard Go error handling often leaves you guessing where an issue actually originated as it bubbles up through intermediate functions, `stacked` captures the context the moment the error is produced.

The linter enforces this "wrap at the source" policy across your entire codebase, offering a seamless debugging experience:

* **Zero Guesswork:** By capturing the exact function, file, and line number where the error occurred, it cuts down debugging time drastically.
* **Foolproof Coverage:** The linter acts as a safety net, guaranteeing that no error is left unwrapped.
* **Frictionless Integration:** Wrapping is idempotent (the first wrap wins) and fully compatible with the standard library (`errors.Is`, `errors.AsType`, and `errors.Unwrap`), meaning your existing error-handling logic remains completely intact.
* **Effortless Adoption:** Migrating an existing codebase doesn't require a tedious manual rewrite. Run the linter with the `-fix` flag to automatically apply wrapping everywhere in a single pass.

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
        slog.Error("panic occurred",
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

## Linter

`stacked-linter` reports every error your code leaves unwrapped.

It only flags errors at the point they cross into your code. That covers:

* Errors returned by functions outside your module (third-party or
  standard-library packages).
* Errors returned by methods called through an interface, since the
  concrete implementation behind it is unknown.
* Errors used from a constant or package-level variable, such as returning
  `sql.ErrNoRows` or your own `ErrNotFound`.
* Errors built from a literal, such as `&MyError{…}`.
* Errors received from channels.

Run it with `-fix` to apply the wrapping automatically, making adoption of `stacked` in an existing codebase a single-pass operation.

There are two ways to run the linter: as a [standalone
binary](#standalone-binary), or as a [golangci-lint
plugin](#golangci-lint-plugin). Both apply
the same rules and accept the same [configuration](#configuration).

### Example

```go
func loadConfig() ([]byte, error) {
	return os.ReadFile("/etc/app/config.yaml") // reported: error returned by os.ReadFile is not wrapped with stacked
}
```

Applying the suggested fix yields:

```go
func loadConfig() ([]byte, error) {
	return stacked.Wrap2(os.ReadFile("/etc/app/config.yaml"))
}
```

### Suppressing a diagnostic

With the standalone binary, use the `//stacked:disable` directive:

```go
err := tx.Rollback() //stacked:disable
```

Under golangci-lint, use the standard `nolint` directive with the linter
name:

```go
err := tx.Rollback() //nolint:stacked
```

### Configuration

Three options tune what the linter considers worth wrapping. The options
are the same whether you run the standalone binary or the golangci-lint
plugin, but the configuration file differs, as shown in each section below.

#### `packages-treated-as-external`

Packages treated as third-party even though they're in your module —
typically generated code. The linter ignores these packages,
but treats errors they return to *your* code as crossing
in from outside, so those calls still need wrapping.

Type: list of package import paths.

```json
["your-module/generated"]
```

#### `ignored-functions`

Functions whose returned error never needs wrapping — typically
error-decorating helpers like `connectrpc.com/connect.NewError` that take an
already-wrapped error and return it, so the trace is already attached.

`errors.AsType`, `errors.Join`, `errors.Unwrap`, are ignored by default.

Type: list of fully-qualified function names, formatted
`<import-path>.<Func>` or `<import-path>.<Type>.<Method>`.

```json
["connectrpc.com/connect.NewError"]
```

#### `check-function-arguments`

Marks a specific function argument as an error supplied for comparison
rather than produced by the program, so the linter leaves that argument
unwrapped. Use it for arguments that receive an existing sentinel error to
check against, like the target argument of `errors.Is`.

Type: list of objects with `function` (a fully-qualified name in the same
format as above) and `argument` (the **1-based** position of the error
argument).

```json
[{ "function": "github.com/stretchr/testify/require.ErrorIs", "argument": 3 }]
```

The `target` arguments of `errors.Is` and `errors.As` are ignored by default.

### Standalone binary

Install and run the `singlechecker` binary:

```sh
go install github.com/tbeati/stacked/linter/cmd/stacked-linter@latest

stacked-linter ./...        # report
stacked-linter -fix ./...   # report and apply suggested fixes
```

[Configuration](#configuration) goes in an optional `stacked.json` in the
working directory:

```json
{
    "packages-treated-as-external": ["example.com/generated"],
    "ignored-functions": ["connectrpc.com/connect.NewError"],
    "check-function-arguments": [
        { "function": "github.com/stretchr/testify/require.ErrorIs", "argument": 3 }
    ]
}
```

### golangci-lint plugin

`stacked` ships as a [golangci-lint module
plugin](https://golangci-lint.run/plugins/module-plugins/). Reference the
plugin in `.custom-gcl.yml`:

```yaml
version: v2.12.2
plugins:
  - module: github.com/tbeati/stacked/linter
    import: github.com/tbeati/stacked/linter/gclplugin
    version: latest
```

Build the custom binary and enable the linter (named `stacked`) in
`.golangci.yml`:

```sh
golangci-lint custom
```

[Configuration](#configuration) goes under the plugin's `settings`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - stacked
  settings:
    custom:
      stacked:
        type: module
        description: Reports errors not wrapped with stacked.
        settings:
          packages-treated-as-external: ["example.com/generated"]
          ignored-functions: ["connectrpc.com/connect.NewError"]
          check-function-arguments:
            - function: github.com/stretchr/testify/require.ErrorIs
              argument: 3
```

Run the resulting `./custom-gcl run ./...` as usual; `--fix` applies the
suggested fixes.