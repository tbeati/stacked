package stacked

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

var (
	ErrNilPanicValue = errors.New("panic with nil value")
	ErrGoexitCalled  = errors.New("runtime.Goexit called")
)

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Error struct {
	Err        error
	StackTrace []StackFrame
}

func (se *Error) Error() string {
	return se.Err.Error()
}

func (se *Error) Unwrap() error {
	return se.Err
}

func Wrap(err error) error {
	return wrap(err, 4)
}

func Wrap2[T any](v T, err error) (T, error) {
	return v, wrap(err, 4)
}

func Wrap3[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, error) {
	return v1, v2, wrap(err, 4)
}

func wrap(err error, skip int) error {
	if err == nil {
		return nil
	}

	var stackErr *Error
	if errors.As(err, &stackErr) {
		return err
	}

	stackTrace := getStackTrace(skip)

	return newError(err, stackTrace)
}

func StackTrace(err error) []StackFrame {
	var stackErr *Error
	if errors.As(err, &stackErr) {
		return stackErr.StackTrace
	}

	return nil
}

func newError(err error, stackTrace []StackFrame) error {
	return &Error{
		Err:        err,
		StackTrace: stackTrace,
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

			err = newError(err, stackTrace)
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
