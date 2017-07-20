package freeportd

import "net"

// GetTCPPort returns either the next available TCP port or an error
// generated during the lookup.
func GetTCPPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
