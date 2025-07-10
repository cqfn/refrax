package protocol

import "net"

// Server defines the interface for a server that can handle incoming A2A messages
type Server interface {
	// Start starts the server and listens on the specified port, while signaling readiness.
	Start(ready chan<- struct{}) error

	// SetHandler sets the message handler for the server.
	SetHandler(handler MessageHandler)

	// Close stops the server gracefully.
	Close() error
}

// MessageHandler defines the function signature for handling messages.
type MessageHandler func(message *Message) (*Message, error)

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
