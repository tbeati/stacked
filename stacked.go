// Package stacked attaches stack traces to errors. [Wrap] captures the call
// stack at its invocation site and stores it alongside the error; [StackTrace]
// retrieves the frames later. Wrapped errors satisfy [errors.Is], [errors.As],
// and [errors.Unwrap], so they compose with standard error handling.
//
// For the captured frames to point at where an error actually originated
// rather than at intermediate forwarders, wrap errors at their source — the
// call or expression that produces them:
//
//	err := stacked.Wrap(os.Chdir("/"))
//
// Wrapping is idempotent: [Wrap] returns nil, already-wrapped errors, and
// ignored errors unchanged, so the first wrap wins. [io.EOF] is ignored by
// default; register additional ignored errors with [Ignore] or [IgnoreFunc].
//
// [Recover] converts panics and [runtime.Goexit] into stacked errors.
package stacked

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"runtime"
)

var ignoredErrors = []error{
	io.EOF,
}

var ignoreFuncs []func(error) bool

// Ignore registers err so future [Wrap], [WrapSeq], or [WrapPull] calls return
// it unchanged instead of attaching a stack trace. Calls are deduplicated by
// identity; registering the same error twice is a no-op. [io.EOF] is registered
// by default.
func Ignore(err error) {
	for _, ignoredError := range ignoredErrors {
		if ignoredError == err {
			return
		}
	}

	ignoredErrors = append(ignoredErrors, err)
}

// IgnoreFunc registers a predicate consulted alongside [Ignore]'s list. If any
// registered predicate returns true for an error, the error is left unwrapped.
func IgnoreFunc(ignoreFunc func(error) bool) {
	ignoreFuncs = append(ignoreFuncs, ignoreFunc)
}

// StackFrame is a single entry in a captured stack trace.
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Error wraps an underlying error together with a stack trace captured at the
// wrap site. Recover one from an arbitrary error via [errors.As] or
// [errors.AsType], or call [StackTrace] for the frames directly.
type Error struct {
	Err        error
	StackTrace []StackFrame
	// iterator marks errors wrapped by WrapSeq/WrapSeq2 so WrapPull can
	// re-wrap with a stack at the pull site rather than at the yield site.
	iterator bool
}

// Error returns the wrapped error's message.
func (se *Error) Error() string {
	return se.Err.Error()
}

// Unwrap returns the wrapped error.
func (se *Error) Unwrap() error {
	return se.Err
}

// Wrap attaches a stack trace to err captured at this call site. It returns
// nil unchanged, ignored errors unchanged, and already-wrapped errors
// unchanged (the first wrap wins). Use the arity-matching variant — [Wrap2],
// [Wrap3], [Wrap4], [Wrap5] — when wrapping the return of a function that
// produces additional values.
func Wrap(err error) error {
	return wrap(err, 4, false)
}

// Wrap2 is the two-value form of [Wrap]; v is returned unchanged.
func Wrap2[T any](v T, err error) (T, error) {
	return v, wrap(err, 4, false)
}

// Wrap3 is the three-value form of [Wrap]; v1 and v2 are returned unchanged.
func Wrap3[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, error) {
	return v1, v2, wrap(err, 4, false)
}

// Wrap4 is the four-value form of [Wrap]; v1, v2, and v3 are returned unchanged.
func Wrap4[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3, error) {
	return v1, v2, v3, wrap(err, 4, false)
}

// Wrap5 is the five-value form of [Wrap]; v1, v2, v3, and v4 are returned unchanged.
func Wrap5[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4, error) {
	return v1, v2, v3, v4, wrap(err, 4, false)
}

// WrapSeq returns an iterator that wraps each error yielded by seq with a
// stack trace captured at yield time.
func WrapSeq(seq iter.Seq[error]) iter.Seq[error] {
	return func(yield func(error) bool) {
		seq(func(err error) bool {
			return yield(wrap(err, 6, true))
		})
	}
}

// WrapSeq2 is the two-value form of [WrapSeq]; values of type T pass through
// unchanged.
func WrapSeq2[T any](seq iter.Seq2[T, error]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		seq(func(v T, err error) bool {
			return yield(v, wrap(err, 6, true))
		})
	}
}

// WrapPull wraps an error returned from an [iter.Pull]-style next function
// with a stack trace captured at the pull site. Already-wrapped errors and
// ignored errors pass through unchanged. The ok value is returned unchanged.
func WrapPull(err error, ok bool) (error, bool) {
	return wrapPull(err), ok
}

// WrapPull2 is the three-value form of [WrapPull]; v is returned unchanged.
func WrapPull2[T any](v T, err error, ok bool) (T, error, bool) {
	return v, wrapPull(err), ok
}

func isIgnored(err error) bool {
	if err == nil {
		return true
	}

	for _, ignoredError := range ignoredErrors {
		if errors.Is(err, ignoredError) {
			return true
		}
	}

	for _, ignoreFunc := range ignoreFuncs {
		if ignoreFunc(err) {
			return true
		}
	}

	return false
}

func wrap(err error, skip int, iterator bool) error {
	if isIgnored(err) {
		return err
	}

	_, ok := errors.AsType[*Error](err)
	if ok {
		return err
	}

	stackTrace := getStackTrace(skip)

	return newError(err, stackTrace, iterator)
}

func wrapPull(err error) error {
	if isIgnored(err) {
		return err
	}

	stackErr, ok := errors.AsType[*Error](err)
	if ok && !stackErr.iterator {
		return err
	}

	stackTrace := getStackTrace(4)

	return newError(err, stackTrace, false)
}

// StackTrace returns the captured frames if err is or wraps a *[Error],
// otherwise nil.
func StackTrace(err error) []StackFrame {
	stackErr, ok := errors.AsType[*Error](err)
	if ok {
		return stackErr.StackTrace
	}

	return nil
}

func newError(err error, stackTrace []StackFrame, iterator bool) error {
	return &Error{
		Err:        err,
		StackTrace: stackTrace,
		iterator:   iterator,
	}
}

func getStackTrace(skip int) []StackFrame {
	const maxFrameCount = 128
	pc := make([]uintptr, maxFrameCount)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]

	var stackTrace []StackFrame
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()

		stackTrace = append(stackTrace, StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return stackTrace
}

var (
	// ErrNilPanicValue is reported by [Recover] when the recovered panic value
	// is nil. On Go ≥ 1.21 this is unreachable without GODEBUG=panicnil=1; the
	// runtime converts panic(nil) into a *[runtime.PanicNilError] instead.
	ErrNilPanicValue = errors.New("panic with nil value")
	// ErrGoexitCalled is reported by [Recover] when its function exits via
	// [runtime.Goexit].
	ErrGoexitCalled = errors.New("runtime.Goexit called")
)

// Recover runs f and routes its outcome through onPanic:
//
//   - f returns normally: onPanic is not called.
//   - f panics: onPanic is called with a *[Error] whose Err is the panic
//     value (if it satisfies error) or [fmt.Errorf]("%v", value) otherwise,
//     and whose stack trace points at the panic site.
//   - f exits via [runtime.Goexit]: onPanic is called with a *[Error]
//     wrapping [ErrGoexitCalled].
//
// onPanic may be nil. If exitOnPanic is true, Recover calls [os.Exit](1)
// after onPanic returns.
func Recover(f func(), onPanic func(err error), exitOnPanic bool) {
	internalRecover(f, func(err error) {
		if onPanic != nil {
			onPanic(err)
		}

		if exitOnPanic {
			os.Exit(1)
		}
	})
}

func internalRecover(f func(), onPanic func(err error)) {
	normalReturn := false
	recovered := false
	var panicValue any
	var stackTrace []StackFrame

	defer func() {
		if !normalReturn {
			var err error
			if recovered {
				if panicValue != nil {
					switch v := panicValue.(type) {
					case error:
						err = v
					default:
						err = fmt.Errorf("%v", v)
					}
				} else {
					err = ErrNilPanicValue
				}
			} else {
				err = ErrGoexitCalled
			}

			err = newError(err, stackTrace, false)
			onPanic(err)
		}
	}()

	func() {
		defer func() {
			panicValue = recover()
			const panicStackSkip = 4
			stackTrace = getStackTrace(panicStackSkip)
		}()

		f()

		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}
