package a

import (
	"net"

	"github.com/tbeati/stacked"
)

func stringLiteralAssignmentExternal() {
	var err error
	_ = err

	err = net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	err = stacked.Wrap(net.UnknownNetworkError("error"))

	_, err = 0, net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	_, err = 0, stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralDeclarationExternal() {
	{
		var err = net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}

	{
		var _, err = 0, net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}
}

func stringLiteralShortDeclarationExternal() {
	{
		err := net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}

	{
		_, err := 0, net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(net.UnknownNetworkError("error"))
		_ = err
	}
}

func stringLiteralReturnSingleExternal() error {
	return net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	return stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralReturnMultipleExternal() (int, error) {
	return 0, net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	return 0, stacked.Wrap(net.UnknownNetworkError("error"))
}

func stringLiteralArgumentExternal() {
	functionWithIntErrorArgument(0, net.UnknownNetworkError("error")) // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(net.UnknownNetworkError("error")))
}

func stringLiteralCompositeLiteralExternal() {
	_ = structWithErrorField{
		err: net.UnknownNetworkError("error"), // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(net.UnknownNetworkError("error")),
	}

	_ = []error{net.UnknownNetworkError("error")} // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	_ = []error{stacked.Wrap(net.UnknownNetworkError("error"))}

	_ = map[string]error{"": net.UnknownNetworkError("error")} // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(net.UnknownNetworkError("error"))}
}

func stringLiteralChannelSendExternal() {
	var errChan chan error

	errChan <- net.UnknownNetworkError("error") // want "value converted to error type net.UnknownNetworkError is not wrapped with stacked"
	errChan <- stacked.Wrap(net.UnknownNetworkError("error"))
}
