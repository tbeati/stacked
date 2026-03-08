package a

import (
	"github.com/tbeati/stacked"
)

func channelReceiveAssignmentInternal() {
	var err error
	_ = err

	err = <-errChan // want "error received from errChan is not wrapped with stacked"
	err = stacked.Wrap(<-errChan)

	_, err = 0, <-errChan // want "error received from errChan is not wrapped with stacked"
	_, err = 0, stacked.Wrap(<-errChan)
}

func channelReceiveDeclarationInternal() {
	{
		var err = <-errChan // want "error received from errChan is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(<-errChan)
		_ = err
	}

	{
		var _, err = 0, <-errChan // want "error received from errChan is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(<-errChan)
		_ = err
	}
}

func channelReceiveShortDeclarationInternal() {
	{
		err := <-errChan // want "error received from errChan is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(<-errChan)
		_ = err
	}

	{
		_, err := 0, <-errChan // want "error received from errChan is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(<-errChan)
		_ = err
	}
}

func channelReceiveReturn1Internal() error {
	return <-errChan // want "error received from errChan is not wrapped with stacked"
	return stacked.Wrap(<-errChan)
}

func channelReceiveReturn2Internal() (int, error) {
	return 0, <-errChan // want "error received from errChan is not wrapped with stacked"
	return 0, stacked.Wrap(<-errChan)
}

func channelReceiveArgumentInternal() {
	functionWithIntErrorArgument(0, <-errChan) // want "error received from errChan is not wrapped with stacked"
	functionWithIntErrorArgument(0, stacked.Wrap(<-errChan))
}

func channelReceiveCompositeLiteralInternal() {
	_ = structWithErrorField{
		err: <-errChan, // want "error received from errChan is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(<-errChan),
	}

	_ = []error{<-errChan} // want "error received from errChan is not wrapped with stacked"
	_ = []error{stacked.Wrap(<-errChan)}

	_ = map[string]error{"": <-errChan} // want "error received from errChan is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(<-errChan)}
}

func channelReceiveChannelSendInternal() {
	var errChan chan error

	errChan <- <-errChan // want "error received from errChan is not wrapped with stacked"
	errChan <- stacked.Wrap(<-errChan)
}
