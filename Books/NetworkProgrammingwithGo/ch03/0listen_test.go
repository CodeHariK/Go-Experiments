// The net.Listen function accepts a network type and an IP address and
// port separated by a colon. The function returns a net.Listener interface
// and an error interface. If the function returns successfully, the listener is
// bound to the specified IP address and port. Binding means that the operating system has exclusively assigned the port on the given IP address to the
// listener. The operating system allows no other processes to listen for incoming traffic on bound ports. If you attempt to bind a listener to a currently
// bound port, net.Listen will return an error.

// You can choose to leave the IP address and port parameters empty. If
// the port is zero or empty, Go will randomly assign a port number to your
// listener.

// You should always be diligent about closing your listener gracefully by
// calling its Close method
// Failure to close the listener
// may lead to memory leaks or deadlocks in your code, because calls to the
// listener’s Accept method may block indefinitely. Closing the listener immediately unblocks calls to the Accept method.

package chapter

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	t.Logf("bound to %q", listener.Addr())
}
