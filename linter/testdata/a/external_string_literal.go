package a

import (
	"net"
)

func stringLiteralAssignmentExternal() {
	var err net.UnknownNetworkError
	_ = err

	err = "error" // want "^basic literal \"error\" is not wrapped with stacked$"

	_, err = 0, "error" // want "^basic literal \"error\" is not wrapped with stacked$"
}

func stringLiteralDeclarationExternal() {
	var err net.UnknownNetworkError = "error" // want "^basic literal \"error\" is not wrapped with stacked$"
	_ = err
}

func stringLiteralReturnSingleExternal() net.UnknownNetworkError {
	return "error" // want "^basic literal \"error\" is not wrapped with stacked$"
}

func stringLiteralReturnMultipleExternal() (int, net.UnknownNetworkError) {
	return 0, "error" // want "^basic literal \"error\" is not wrapped with stacked$"
}

func stringLiteralArgumentExternal() {
	functionWithIntStringTypeErrorArgument(0, "error") // want "^basic literal \"error\" is not wrapped with stacked$"
}

func stringLiteralCompositeLiteralExternal() {
	_ = structWithStringTypeErrorField{
		err: "error", // want "^basic literal \"error\" is not wrapped with stacked$"
	}

	_ = []net.UnknownNetworkError{"error"} // want "^basic literal \"error\" is not wrapped with stacked$"

	_ = map[string]net.UnknownNetworkError{"": "error"} // want "^basic literal \"error\" is not wrapped with stacked$"
}

func stringLiteralChannelSendExternal() {
	var errChan chan net.UnknownNetworkError

	errChan <- "error" // want "^basic literal \"error\" is not wrapped with stacked$"
}
