package server

import (
	"context"
	"net"
)

// Server interface defines requirements for servers
type Server interface {
	Listen(context.Context, string) (net.Listener, error)
	WaitCommands()
}
