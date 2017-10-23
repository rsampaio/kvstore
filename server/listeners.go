package server

import (
	"crypto/tls"
	"net"
)

// NewTCPListener creates a TCP listener.
func NewTCPListener(hostport string) (net.Listener, error) {
	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		return nil, err
	}
	return ln, err
}

// NewTLSListener creates a TLS listener.
func NewTLSListener(hostport string, config *tls.Config) (net.Listener, error) {
	ln, err := tls.Listen("tcp", hostport, config)
	if err != nil {
		return nil, err
	}
	return ln, err
}
