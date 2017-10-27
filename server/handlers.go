package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/rsampaio/kvstore/store"
)

// HandlerFunc define the function to handle each command
type HandlerFunc func(store.Store, []string, *bufio.Reader, net.Conn) (string, error)

// Handlers is a map of HandlerFunc using commands as keys
type Handlers map[string]HandlerFunc

// Handler implements each operation supported by the protocol.
type Handler struct{}

var defaultHandler = Handler{}

// DefaultHandlers define set functions to handle commands.
var DefaultHandlers = Handlers{
	"SET":    defaultHandler.Set,
	"GET":    defaultHandler.Get,
	"DELETE": defaultHandler.Delete,
	"STREAM": defaultHandler.Stream,
}

// Set receives a store, slice of args and a value to store
// and returns a reply and an error.
func (h Handler) Set(s store.Store, args []string, in *bufio.Reader, conn net.Conn) (string, error) {
	key := args[0]
	size, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return "ERROR\r", err
	}

	buf := bytes.NewBufferString("")
	io.CopyN(buf, in, size)
	return "OK\r", s.Set(key, buf.String())
}

// Get receives a store, a slice of args and a connection and
// handles the GET command when it is parsed by the protocol.
func (h Handler) Get(s store.Store, args []string, _ *bufio.Reader, conn net.Conn) (string, error) {
	value, _ := s.Get(args[0])
	fmt.Fprintf(conn, "VALUE %d\r\n", len(value))
	buf := bytes.NewReader([]byte(value))
	io.CopyN(conn, buf, int64(len(value)))
	return "\r", nil
}

// Delete receives a store, a slice of args and a connection and
// handles the GET command when it is parsed by the protocol.
func (h Handler) Delete(s store.Store, args []string, _ *bufio.Reader, _ net.Conn) (string, error) {
	return "OK\r", s.Delete(args[0])
}

// Stream sends all keys with associated values ordered by last modified time
func (h Handler) Stream(s store.Store, _ []string, _ *bufio.Reader, conn net.Conn) (string, error) {
	list := s.GetLastModifiedKeys()
	for _, k := range list {
		v, _ := s.Get(k)
		fmt.Fprintln(conn, fmt.Sprintf("%v %v\r", k, v))
	}
	return "OK\r", nil
}
