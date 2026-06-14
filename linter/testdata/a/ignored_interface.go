package a

import (
	"testdata/b"
)

func ignoredInterfaceMethodCallAssignment() {
	var err error
	_ = err
	var i b.IgnoredInterface

	err = i.SingleReturn()
	_, err = 0, i.SingleReturn()
	_, err = i.MultipleReturn()
}

func ignoredInterfaceMethodCallDeclaration() {
	var i b.IgnoredInterface

	var err = i.SingleReturn()
	_ = err
	var _, err2 = i.MultipleReturn()
	_ = err2
}

func ignoredInterfaceMethodCallShortDeclaration() {
	var i b.IgnoredInterface

	err := i.SingleReturn()
	_ = err
	_, err2 := i.MultipleReturn()
	_ = err2
}

func ignoredInterfaceMethodCallReturn1() error {
	var i b.IgnoredInterface

	return i.SingleReturn()
}

func ignoredInterfaceMethodCallReturn2() (int, error) {
	var i b.IgnoredInterface

	return 0, i.SingleReturn()
	return i.MultipleReturn()
}

func ignoredInterfaceMethodCallArgument() {
	var i b.IgnoredInterface

	functionWithErrorArgument(i.SingleReturn())
	functionWithIntErrorArgument(0, i.SingleReturn())
	functionWithIntErrorArgument(i.MultipleReturn())
}

func ignoredInterfaceMethodCallCompositeLiteral() {
	var i b.IgnoredInterface

	_ = structWithErrorField{
		err: i.SingleReturn(),
	}
	_ = []error{i.SingleReturn()}
	_ = map[string]error{"": i.SingleReturn()}
}

func ignoredInterfaceMethodCallChannelSend() {
	var errChan chan error
	var i b.IgnoredInterface

	errChan <- i.SingleReturn()
}
