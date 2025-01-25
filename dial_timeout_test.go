package main

import (
	"errors"
	"net"
	"syscall"
	"testing"
	"time"
)

func DailTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, addr)
}

func TestDialTimeout(t *testing.T) {
	conn, err := DailTimeout("tcp", "127.0.0.1:8080", 5*time.Second)
	if err == nil {
		err := conn.Close()
		if err != nil {
			return
		}
		t.Fatalf("connection did not timeout")
	}
	var netErr net.Error
	ok := errors.As(err, &netErr)
	if !ok {
		t.Fatalf("error does not implement net.Error: %v", err)
	}
	if !netErr.Timeout() {
		t.Fatalf("error is not a timeout")
	}
}
