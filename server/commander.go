package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rsampaio/kvstore/protocol"
	"github.com/rsampaio/kvstore/store"
)

type internalMetrics struct {
	clientCount int
}

// Commander has a store.Store field that is passed to default handlers
type Commander struct {
	store    store.Store
	listener net.Listener
	metrics  internalMetrics
}

// NewCommander receives a store and a listener and returns a new Commander instance
func NewCommander(store store.Store, list net.Listener) *Commander {
	return &Commander{
		store:    store,
		listener: list,
	}
}

// Run runs a loop accepting connections to the listener and executes the commander
func (c *Commander) Run(ctx context.Context) error {
	defer c.listener.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				fmt.Printf(
					"listener=%v capacity-left=%dbytes\n",
					c.listener.Addr(),
					c.store.Cap(),
				)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			break
		default:
			conn, err := c.listener.Accept()
			if err != nil {
				return err
			}
			go func() {
				if err := c.WaitCommands(ctx, conn); err != nil {
					fmt.Printf("wait command error: %v\n", err)
					conn.Close()
				}
			}()
		}
	}
}

// WaitCommands handles connections and parse commands from clients
func (c *Commander) WaitCommands(ctx context.Context, conn net.Conn) error {
	p := &protocol.Protocol{}

	buf := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			break
		default:
			// Ignore empty lines
			line, _, err := buf.ReadLine()
			if string(line) == "" {
				continue
			}

			// If command is invalid parse will return an error
			if err := p.Parse(string(line)); err != nil {
				fmt.Fprintln(conn, "ERROR")
				return err
			}

			result, err := DefaultHandlers[p.Command](c.store, p.Args, buf, conn)

			// Adds \n back but not \r
			fmt.Fprintln(conn, result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
