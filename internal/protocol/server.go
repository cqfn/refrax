package protocol

import "net"

type Server interface {
	Start(ready chan<- struct{}) error
	SetHandler(handler MessageHandler)
	Close() error
}

type MessageHandler func(message *Message) (*Message, error)

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
