package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/rsampaio/kvstore/server"
	"github.com/rsampaio/kvstore/store"
)

var (
	enableTLS = flag.Bool("enable-tls", false, "Enables TLS server (requires --tls-cert and --tls-key)")
	tcpPort   = flag.String("tcp-listen", ":2020", "TCP server listen address")
	tlsPort   = flag.String("tls-listen", ":2021", "TLS server listen address")
	tlsCert   = flag.String("tls-cert", "", "PEM certificate file")
	tlsKey    = flag.String("tls-key", "", "Cerficate key file")
	capacity  = flag.Int("capacity-bytes", 1000, "Max capacity in bytes")
)

func startTCP(ctx context.Context, s store.Store) {
	fmt.Printf("starting-tcp port=%v\n", *tcpPort)
	l, err := server.NewTCPListener(*tcpPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	r := server.NewCommander(s, l)
	go func(ctx context.Context) {
		r.Run(ctx)
	}(ctx)
}

func startTLS(ctx context.Context, s store.Store) {
	tlsPort := *tlsPort
	if *tlsCert == "" || *tlsKey == "" {
		fmt.Fprintln(os.Stderr, "missing --tls-cert or --tls-key arguments")
		return
	}

	fmt.Printf("starting-tls cert=%v key=%v port=%v\n", *tlsCert, *tlsKey, tlsPort)

	cert, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading certificates %v\n", err.Error())
		return
	}

	ls, err := server.NewTLSListener(tlsPort, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	rs := server.NewCommander(s, ls)

	go func() {
		rs.Run(ctx)
	}()
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel context on SIGQUIT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	s := store.NewMemoryStore(*capacity)

	startTCP(ctx, s)
	if *enableTLS {
		startTLS(ctx, s)
	}

	select {
	case <-ctx.Done():
		return
	case <-sigCh:
		fmt.Printf("shutting-down\n")
		cancel()
	}
}
