// Package server defines the types required for the server workflow.
//
// The server also defines a Commander type that execute commands
// after each command line is parsed by the protocol, a set of default handlers
// that respond to each of the commands defined by the protocol as well as
// helper functions to initalize a TCP and a TLS listener.
package server
