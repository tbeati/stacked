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

func Ignore(err error) {
	for _, ignoredError := range ignoredErrors {
		if ignoredError == err {
			return
		}
	}

	ignoredErrors = append(ignoredErrors, err)
}

func IgnoreFunc(ignoreFunc func(error) bool) {
	ignoreFuncs = append(ignoreFuncs, ignoreFunc)
}

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Error struct {
	Err        error
	StackTrace []StackFrame
	iterator   bool
}

func (se *Error) Error() string {
	return se.Err.Error()
}

func (se *Error) Unwrap() error {
	return se.Err
}

func Wrap(err error) error {
	return wrap(err, 4, false)
}

func Wrap2[T any](v T, err error) (T, error) {
	return v, wrap(err, 4, false)
}

func Wrap3[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, error) {
	return v1, v2, wrap(err, 4, false)
}

func Wrap4[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3, error) {
	return v1, v2, v3, wrap(err, 4, false)
}

func Wrap5[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4, error) {
	return v1, v2, v3, v4, wrap(err, 4, false)
}

func WrapSeq(seq iter.Seq[error]) iter.Seq[error] {
	return func(yield func(error) bool) {
		seq(func(err error) bool {
			return yield(wrap(err, 6, true))
		})
	}
}

func WrapSeq2[T any](seq iter.Seq2[T, error]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		seq(func(v T, err error) bool {
			return yield(v, wrap(err, 6, true))
		})
	}
}

func WrapPull(err error, ok bool) (error, bool) {
	return wrapPull(err), ok
}

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
	ErrNilPanicValue = errors.New("panic with nil value")
	ErrGoexitCalled  = errors.New("runtime.Goexit called")
)

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
