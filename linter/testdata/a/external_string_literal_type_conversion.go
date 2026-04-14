package a

import (
	"net"

	"github.com/tbeati/stacked"
)

func stringLiteralTypeConversionAssignmentExternal() {
	var err error
	_ = err

	err = net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	err = stacked.Wrap(net.UnknownNetworkError("error"))

	_, err = 0, net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	_, err = 0, stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralTypeConversionDeclarationExternal() {
	{
		var err = net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
		_ = err
	}
	{
		var err = stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}

	{
		var _, err = 0, net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}
}

func stringLiteralTypeConversionShortDeclarationExternal() {
	{
		err := net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
		_ = err
	}
	{
		err := stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}

	{
		_, err := 0, net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}
}

func stringLiteralTypeConversionReturnSingleExternal() error {
	return net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	return stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralTypeConversionReturnMultipleExternal() (int, error) {
	return 0, net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	return 0, stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralTypeConversionArgumentExternal() {
	functionWithIntErrorArgument(0, net.UnknownNetworkError("error")) // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	functionWithIntErrorArgument(0, stacked.Wrap(net.UnknownNetworkError("error")))
}

func stringLiteralTypeConversionCompositeLiteralExternal() {
	_ = structWithErrorField{
		err: net.UnknownNetworkError("error"), // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(net.UnknownNetworkError("error")),
	}

	_ = []error{net.UnknownNetworkError("error")} // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	_ = []error{stacked.Wrap(net.UnknownNetworkError("error"))}

	_ = map[string]error{"": net.UnknownNetworkError("error")} // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	_ = map[string]error{"": stacked.Wrap(net.UnknownNetworkError("error"))}
}

func stringLiteralTypeConversionChannelSendExternal() {
	var errChan chan error

	errChan <- net.UnknownNetworkError("error") // want "^value converted to error type net.UnknownNetworkError is not wrapped with stacked$"
	errChan <- stacked.Wrap(net.UnknownNetworkError("error"))
}
