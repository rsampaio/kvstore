package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/rsampaio/kvstore/store"
)

func TestServer(t *testing.T) {
	st := store.NewMemoryStore(100)
	ln, _ := NewTCPListener("localhost:10000")
	s := NewCommander(st, ln)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := s.Run(ctx); err != nil {
			t.Logf("server error: %v", err)
		}
	}()

	t.Run("client", func(t *testing.T) {
		c, err := net.Dial("tcp", "localhost:10000")
		defer c.Close()

		if err != nil {
			t.Errorf("unexpected connect error: %v", err)
		}

		if _, err := fmt.Fprint(c, "SET foo 1\r\n"); err != nil {
			t.Fatalf("unexpected client write error: %v", err)
		}

		if _, err = fmt.Fprint(c, "a"); err != nil {
			t.Fatalf("unexpected write error: %v", err)
		}

		r, _ := bufio.NewReader(c).ReadString('\n')
		if r != "OK\r\n" {
			t.Fatalf("unexpected SET response: %v", r)
		}
		c.Close()
	})
}

func BenchmarkServer(b *testing.B) {
	st := store.NewMemoryStore(100)
	ln, _ := NewTCPListener("localhost:10001")
	defer ln.Close()

	s := NewCommander(st, ln)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s.Run(ctx)
	}()

	c, _ := net.Dial("tcp", "localhost:10001")
	defer c.Close()

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			fmt.Fprint(c, "SET foo 1\r\n")
			fmt.Fprint(c, "a")

			v, err := bufio.NewReader(c).ReadString('\n')
			if err != nil && v != "OK\r\n" {
				b.Errorf("unexpected response")
			}
		}
	})

	b.ResetTimer()
	b.Run("Get", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			fmt.Fprint(c, "GET foo\r\n")
			buf := bufio.NewReader(c)
			v, err := buf.ReadString('\n')
			if err != nil && v != "VALUE 1" {
				b.Errorf("unexpected response: %v", v)
			}

			v, err = buf.ReadString('\n')
			if err != nil && v != "a\r\n" {
				b.Errorf("unexpected response: %v", v)
			}
		}
	})
	cancel()
}
