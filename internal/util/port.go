package util

import "net"

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
