package stacked

import (
	"errors"
	"io"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
)

// saveIgnoreState snapshots the package-level ignore state and restores it
// when the test ends. Required because Ignore / IgnoreFunc mutate globals.
func saveIgnoreState(t *testing.T) {
	t.Helper()
	origErrors := slices.Clone(ignoredErrors)
	origFuncs := slices.Clone(ignoreFuncs)
	t.Cleanup(func() {
		ignoredErrors = origErrors
		ignoreFuncs = origFuncs
	})
}

func TestWrap_Nil(t *testing.T) {
	err := Wrap(nil)
	if err != nil {
		t.Errorf("Wrap(nil) = %v; want nil", err)
	}
}

func TestWrap_IgnoredError(t *testing.T) {
	err := Wrap(io.EOF)
	if err != io.EOF { //nolint:errorlint
		t.Errorf("Wrap(io.EOF) = %v; want io.EOF unchanged", err)
	}
}

func TestWrap_FreshError(t *testing.T) {
	orig := errors.New("boom")
	wrapped := Wrap(orig)

	se, ok := errors.AsType[*Error](wrapped)
	if !ok {
		t.Fatalf("Wrap(orig) is %T; want *Error", wrapped)
	}
	if se.Err != orig { //nolint:errorlint
		t.Errorf("wrapped.Err = %v; want %v", se.Err, orig)
	}
	if len(se.StackTrace) == 0 {
		t.Fatal("StackTrace is empty")
	}
	if !strings.Contains(se.StackTrace[0].Function, "TestWrap_FreshError") {
		t.Errorf("top frame = %q; want it to mention TestWrap_FreshError", se.StackTrace[0].Function)
	}
}

func TestWrap_DoubleWrapReturnsSame(t *testing.T) {
	orig := errors.New("boom")
	once := Wrap(orig)
	twice := Wrap(once)
	if once != twice { //nolint:errorlint
		t.Error("Wrap(Wrap(e)) produced a new error; want same pointer")
	}
}

func TestWrap2_PassesValuesThrough(t *testing.T) {
	orig := errors.New("boom")
	v, err := Wrap2(42, orig)
	if v != 42 {
		t.Errorf("v = %d; want 42", v)
	}
	_, ok := errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestWrap3_PassesValuesThrough(t *testing.T) {
	orig := errors.New("boom")
	a, b, err := Wrap3(1, "two", orig)
	if a != 1 || b != "two" {
		t.Errorf("values = %d, %q; want 1, \"two\"", a, b)
	}
	_, ok := errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestWrap4_PassesValuesThrough(t *testing.T) {
	orig := errors.New("boom")
	a, b, c, err := Wrap4(1, "two", 3.5, orig)
	if a != 1 || b != "two" || c != 3.5 {
		t.Errorf("values = %d, %q, %v; want 1, \"two\", 3.5", a, b, c)
	}
	_, ok := errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestWrap5_PassesValuesThrough(t *testing.T) {
	orig := errors.New("boom")
	a, b, c, d, err := Wrap5(1, "two", 3.5, true, orig)
	if a != 1 || b != "two" || c != 3.5 || d != true { //nolint:revive
		t.Errorf("values = %d, %q, %v, %v", a, b, c, d)
	}
	_, ok := errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestIgnore_SkipsWrap(t *testing.T) {
	saveIgnoreState(t)
	sentinel := errors.New("ignore me")
	Ignore(sentinel)
	got := Wrap(sentinel)
	if got != sentinel { //nolint:errorlint
		t.Errorf("Wrap(ignored) = %v; want sentinel unchanged", got)
	}
}

func TestIgnore_Deduplicates(t *testing.T) {
	saveIgnoreState(t)
	sentinel := errors.New("dup")
	before := len(ignoredErrors)
	Ignore(sentinel)
	Ignore(sentinel)
	added := len(ignoredErrors) - before
	if added != 1 {
		t.Errorf("added %d entries for two Ignore calls; want 1", added)
	}
}

func TestIgnoreFunc(t *testing.T) {
	saveIgnoreState(t)
	IgnoreFunc(func(err error) bool {
		return err.Error() == "ignore-by-func"
	})
	err := errors.New("ignore-by-func")
	got := Wrap(err)
	if got != err { //nolint:errorlint
		t.Errorf("Wrap(matched-by-IgnoreFunc) = %v; want unchanged", got)
	}
}

func TestWrapSeq(t *testing.T) {
	orig := errors.New("boom")
	seq := func(yield func(error) bool) {
		yield(orig)
		yield(orig)
	}

	var collected []error
	for err := range WrapSeq(seq) {
		collected = append(collected, err)
	}
	if len(collected) != 2 {
		t.Fatalf("yielded %d errors; want 2", len(collected))
	}
	for i, err := range collected {
		se, ok := errors.AsType[*Error](err)
		if !ok {
			t.Fatalf("collected[%d] is %T; want *Error", i, err)
		}
		if !se.iterator {
			t.Errorf("collected[%d].iterator = false; want true", i)
		}
		if !strings.Contains(se.StackTrace[0].Function, "TestWrapSeq") {
			t.Errorf("collected[%d] top frame = %q; want TestWrapSeq", i, se.StackTrace[0].Function)
		}
	}
}

func TestWrapSeq2(t *testing.T) {
	orig := errors.New("boom")
	seq := func(yield func(int, error) bool) {
		yield(1, orig)
	}

	var gotV int
	var gotErr error
	for v, err := range WrapSeq2(seq) {
		gotV = v
		gotErr = err
	}
	if gotV != 1 {
		t.Errorf("v = %d; want 1", gotV)
	}
	se, ok := errors.AsType[*Error](gotErr)
	if !ok {
		t.Fatalf("err is %T; want *Error", gotErr)
	}
	if !se.iterator {
		t.Error("iterator flag = false; want true")
	}
}

func TestWrapPull_PlainError(t *testing.T) {
	orig := errors.New("boom")
	err, ok := WrapPull(orig, true)
	if !ok {
		t.Error("ok = false; want true (pass-through)")
	}
	_, ok = errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestWrapPull_AlreadyWrappedNonIterator(t *testing.T) {
	orig := errors.New("boom")
	wrapped := Wrap(orig)
	out, _ := WrapPull(wrapped, true)
	if out != wrapped { //nolint:errorlint
		t.Error("WrapPull rewrapped a non-iterator *Error; want pass-through")
	}
}

func TestWrapPull_AfterWrapSeq(t *testing.T) {
	orig := errors.New("boom")
	seq := func(yield func(error) bool) {
		yield(orig)
	}

	var fromSeq error
	for err := range WrapSeq(seq) {
		fromSeq = err
	}

	out, _ := WrapPull(fromSeq, true)
	if out == fromSeq { //nolint:errorlint
		t.Fatal("WrapPull did not rewrap iterator-flagged error")
	}
	se, ok := errors.AsType[*Error](out)
	if !ok {
		t.Fatalf("out is %T; want *Error", out)
	}
	if se.iterator {
		t.Error("rewrapped error still has iterator=true")
	}
}

func TestWrapPull_IgnoredError(t *testing.T) {
	out, ok := WrapPull(io.EOF, true)
	if out != io.EOF { //nolint:errorlint
		t.Errorf("WrapPull(io.EOF) = %v; want io.EOF unchanged", out)
	}
	if !ok {
		t.Error("ok = false; want true (pass-through)")
	}
}

func TestWrapPull2(t *testing.T) {
	orig := errors.New("boom")
	v, err, ok := WrapPull2(42, orig, true)
	if v != 42 || !ok {
		t.Errorf("pass-through dropped: v=%d ok=%t; want 42, true", v, ok)
	}
	_, ok = errors.AsType[*Error](err) //nolint:errcheck
	if !ok {
		t.Fatalf("err is %T; want *Error", err)
	}
}

func TestStackTrace_Wrapped(t *testing.T) {
	err := Wrap(errors.New("boom"))
	frames := StackTrace(err)
	if len(frames) == 0 {
		t.Fatal("StackTrace returned no frames for wrapped error")
	}
}

func TestStackTrace_Plain(t *testing.T) {
	frames := StackTrace(errors.New("boom"))
	if frames != nil {
		t.Errorf("StackTrace(plain) = %v; want nil", frames)
	}
}

func TestError_Methods(t *testing.T) {
	orig := errors.New("boom")
	wrapped := Wrap(orig).(*Error) //nolint:errorlint
	if wrapped.Error() != "boom" {
		t.Errorf("Error() = %q; want \"boom\"", wrapped.Error())
	}
	if wrapped.Unwrap() != orig { //nolint:errorlint
		t.Errorf("Unwrap() = %v; want orig", wrapped.Unwrap())
	}
}

func TestRecover_PanicString(t *testing.T) {
	var captured error
	Recover(func() {
		panic("oops")
	}, func(err error) {
		captured = err
	}, false)

	if captured == nil {
		t.Fatal("onPanic not called")
	}
	_, ok := errors.AsType[*Error](captured) //nolint:errcheck
	if !ok {
		t.Fatalf("captured is %T; want *Error", captured)
	}
	if !strings.Contains(captured.Error(), "oops") {
		t.Errorf("error message = %q; want it to contain \"oops\"", captured.Error())
	}
}

func TestRecover_PanicError(t *testing.T) {
	orig := errors.New("boom")
	var captured error
	Recover(func() {
		panic(orig)
	}, func(err error) {
		captured = err
	}, false)

	if !errors.Is(captured, orig) {
		t.Errorf("captured = %v; want errors.Is(_, orig)", captured)
	}
	frames := StackTrace(captured)
	if len(frames) == 0 {
		t.Error("captured panic has no stack trace")
	}
}

func TestRecover_Goexit(t *testing.T) {
	var captured error
	var wg sync.WaitGroup
	wg.Go(func() {
		Recover(func() {
			runtime.Goexit()
		}, func(err error) {
			captured = err
		}, false)
	})
	wg.Wait()

	if !errors.Is(captured, ErrGoexitCalled) {
		t.Errorf("captured = %v; want ErrGoexitCalled", captured)
	}
}

func TestRecover_NormalReturnSkipsOnPanic(t *testing.T) {
	called := false
	Recover(func() {
		// normal return
	}, func(err error) {
		called = true
	}, false)

	if called {
		t.Error("onPanic invoked for non-panicking function")
	}
}

func TestRecover_NilOnPanicSafe(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Errorf("Recover re-panicked: %v", r)
		}
	}()
	Recover(func() {
		panic("oops")
	}, nil, false)
}
