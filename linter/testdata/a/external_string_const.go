package a

import (
	"net"
)

const errMessage = "error"

func stringConstAssignmentExternal() {
	var err net.UnknownNetworkError
	_ = err

	err = errMessage // want "^errMessage is not wrapped with stacked$"

	_, err = 0, errMessage // want "^errMessage is not wrapped with stacked$"
}

func stringConstDeclarationExternal() {
	var err net.UnknownNetworkError = errMessage // want "^errMessage is not wrapped with stacked$"
	_ = err
}

func stringConstReturnSingleExternal() net.UnknownNetworkError {
	return errMessage // want "^errMessage is not wrapped with stacked$"
}

func stringConstReturnMultipleExternal() (int, net.UnknownNetworkError) {
	return 0, errMessage // want "^errMessage is not wrapped with stacked$"
}

func stringConstArgumentExternal() {
	functionWithIntStringTypeErrorArgument(0, errMessage) // want "^errMessage is not wrapped with stacked$"
}

func stringConstCompositeLiteralExternal() {
	_ = structWithStringTypeErrorField{
		err: errMessage, // want "^errMessage is not wrapped with stacked$"
	}

	_ = []net.UnknownNetworkError{errMessage} // want "^errMessage is not wrapped with stacked$"

	_ = map[string]net.UnknownNetworkError{"": errMessage} // want "^errMessage is not wrapped with stacked$"
}

func stringConstChannelSendExternal() {
	var errChan chan net.UnknownNetworkError

	errChan <- errMessage // want "^errMessage is not wrapped with stacked$"
}
