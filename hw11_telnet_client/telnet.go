package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telenet struct {
	Address string
	Timeout time.Duration
	Conn    net.Conn
	In      io.ReadCloser
	Out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telenet{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

func (t *Telenet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	if err != nil {
		return err
	}

	t.Conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.Address)
	return nil
}

func (t *Telenet) Close() error {
	return t.Conn.Close()
}

func (t *Telenet) Send() error {
	_, err := io.Copy(t.Conn, t.In)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (t *Telenet) Receive() error {
	_, err := io.Copy(t.Out, t.Conn)
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}
