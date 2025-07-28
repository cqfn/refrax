package protocol

import (
	"context"
	"net"
)

// Server defines the interface for a server that can handle incoming A2A messages
type Server interface {
	// ListenAndServe starts the server and listens on the specified port, while signaling readiness.
	ListenAndServe() error

	// MsgHandler sets the message handler for the server.
	MsgHandler(handler MsgHandler)

	// Handler sets the handler function for processing requests.
	Handler(handler Handler)

	// Shutdown stops the server gracefully.
	Shutdown() error

	// Ready returns a channel that signals when the server is ready to accept requests.
	Ready() <-chan bool
}

type (
	// Handler handles all incoming requests on JSONRPC level
	Handler func(next Handler, r *JSONRPCRequest) (*JSONRPCResponse, error)
	// MsgHandler handles messages received from the A2A server on Message level
	MsgHandler func(ctx context.Context, message *Message) (*Message, error)
)

// FreePort finds a free TCP port on the localhost.
func FreePort() (port int, err error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := l.Close(); cerr != nil {
			err = cerr
		}
	}()
	port = l.Addr().(*net.TCPAddr).Port
	return port, nil
}
