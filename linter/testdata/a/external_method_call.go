package a

import (
	"net"
	"net/netip"
	"os"

	"github.com/tbeati/stacked"
)

func methodCallAssignmentExternal() {
	var err error
	_ = err
	var conn *net.UDPConn

	err = conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
	err = stacked.Wrap(conn.Close())

	_, err = 0, conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
	_, err = 0, stacked.Wrap(conn.Close())

	_, err = conn.File() // want "error returned by conn.File is not wrapped with stacked"
	_, err = stacked.Wrap2(conn.File())

	_, _, _, _, err = conn.ReadMsgUDPAddrPort(nil, nil) // want "error returned by conn.ReadMsgUDPAddrPort is not wrapped with stacked"
	_, _, _, _, err = stacked.Wrap5(conn.ReadMsgUDPAddrPort(nil, nil))
}

func methodCallDeclarationExternal() {
	var conn *net.UDPConn

	{
		var err = conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
		_ = err
	}
	{
		var err = stacked.Wrap(conn.Close())
		_ = err
	}

	{
		var _, err = 0, conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = 0, stacked.Wrap(conn.Close())
		_ = err
	}

	{
		var _, err = conn.File() // want "error returned by conn.File is not wrapped with stacked"
		_ = err
	}
	{
		var _, err = stacked.Wrap2(conn.File())
		_ = err
	}

	{
		var _, _, _, _, err = conn.ReadMsgUDPAddrPort(nil, nil) // want "error returned by conn.ReadMsgUDPAddrPort is not wrapped with stacked"
		_ = err
	}
	{
		var _, _, _, _, err = stacked.Wrap5(conn.ReadMsgUDPAddrPort(nil, nil))
		_ = err
	}
}

func methodCallShortDeclarationExternal() {
	var conn *net.UDPConn

	{
		err := conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
		_ = err
	}
	{
		err := stacked.Wrap(conn.Close())
		_ = err
	}

	{
		_, err := 0, conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
		_ = err
	}
	{
		_, err := 0, stacked.Wrap(conn.Close())
		_ = err
	}

	{
		_, err := conn.File() // want "error returned by conn.File is not wrapped with stacked"
		_ = err
	}
	{
		_, err := stacked.Wrap2(conn.File())
		_ = err
	}

	{
		_, _, _, _, err := conn.ReadMsgUDPAddrPort(nil, nil) // want "error returned by conn.ReadMsgUDPAddrPort is not wrapped with stacked"
		_ = err
	}
	{
		_, _, _, _, err := stacked.Wrap5(conn.ReadMsgUDPAddrPort(nil, nil))
		_ = err
	}
}

func methodCallReturn1External() error {
	var conn *net.UDPConn

	return conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
	return stacked.Wrap(conn.Close())
}

func methodCallReturn2External() (*os.File, error) {
	var conn *net.UDPConn

	return nil, conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
	return nil, stacked.Wrap(conn.Close())

	return conn.File() // want "error returned by conn.File is not wrapped with stacked"
	return stacked.Wrap2(conn.File())
}

func methodCallReturn5External() (int, int, int, netip.AddrPort, error) {
	var conn *net.UDPConn

	return 0, 0, 0, netip.AddrPort{}, conn.Close() // want "error returned by conn.Close is not wrapped with stacked"
	return 0, 0, 0, netip.AddrPort{}, stacked.Wrap(conn.Close())

	return conn.ReadMsgUDPAddrPort(nil, nil) // want "error returned by conn.ReadMsgUDPAddrPort is not wrapped with stacked"
	return stacked.Wrap5(conn.ReadMsgUDPAddrPort(nil, nil))
}

func methodCallArgumentExternal() {
	var conn *net.UDPConn

	functionWithErrorArgument(conn.Close()) // want "error returned by conn.Close is not wrapped with stacked"
	functionWithErrorArgument(stacked.Wrap(conn.Close()))

	functionWithFileErrorArgument(nil, conn.Close()) // want "error returned by conn.Close is not wrapped with stacked"
	functionWithFileErrorArgument(nil, stacked.Wrap(conn.Close()))

	functionWithFileErrorArgument(conn.File()) // want "error returned by conn.File is not wrapped with stacked"
	functionWithFileErrorArgument(stacked.Wrap2(conn.File()))

	functionWithIntIntIntAddrPortErrorArgument(0, 0, 0, netip.AddrPort{}, conn.Close()) // want "error returned by conn.Close is not wrapped with stacked"
	functionWithIntIntIntAddrPortErrorArgument(0, 0, 0, netip.AddrPort{}, stacked.Wrap(conn.Close()))

	functionWithIntIntIntAddrPortErrorArgument(conn.ReadMsgUDPAddrPort(nil, nil)) // want "error returned by conn.ReadMsgUDPAddrPort is not wrapped with stacked"
	functionWithIntIntIntAddrPortErrorArgument(stacked.Wrap5(conn.ReadMsgUDPAddrPort(nil, nil)))
}

func methodCallCompositeLiteralExternal() {
	var file *os.File

	_ = structWithErrorField{
		err: file.Chdir(), // want "error returned by file.Chdir is not wrapped with stacked"
	}
	_ = structWithErrorField{
		err: stacked.Wrap(file.Chdir()),
	}

	_ = []error{file.Chdir()} // want "error returned by file.Chdir is not wrapped with stacked"
	_ = []error{stacked.Wrap(file.Chdir())}

	_ = map[string]error{"": file.Chdir()} // want "error returned by file.Chdir is not wrapped with stacked"
	_ = map[string]error{"": stacked.Wrap(file.Chdir())}
}

func methodCallChannelSendExternal() {
	var errChan chan error
	var file *os.File

	errChan <- file.Chdir() // want "error returned by file.Chdir is not wrapped with stacked"
	errChan <- stacked.Wrap(file.Chdir())
}

func methodCallBlankAssignmentExternal() {
	var conn *net.UDPConn

	_ = conn.Close()
	_, _ = 0, conn.Close()
	_, _ = conn.File()
	_, _, _, _, _ = conn.ReadMsgUDPAddrPort(nil, nil)
}
