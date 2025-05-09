package main

import (
	"bytes"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestOK(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	payload := "hello, kjdnawakwjdnawkdjawjnd!\n"

	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)
		defer conn.Close()

		buf := make([]byte, len(payload))
		n, err := io.ReadFull(conn, buf)
		require.NoError(t, err)
		require.Equal(t, len(payload), n)
		require.Equal(t, payload, string(buf))

		n, err = conn.Write(buf)
		require.NoError(t, err)
		require.Equal(t, len(buf), n)
	}()

	in := io.NopCloser(strings.NewReader(payload))
	out := &bytes.Buffer{}

	client := NewTelnetClient(l.Addr().String(), time.Second, in, out)
	require.NoError(t, client.Connect())
	defer client.Close()

	require.NoError(t, client.Send())
	require.NoError(t, client.Receive())
	require.Equal(t, payload, out.String())
}

func TestConnectTimeout(t *testing.T) {
	client := NewTelnetClient("google.com:65000", 100*time.Millisecond, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
	start := time.Now()
	err := client.Connect()
	require.Error(t, err)
	require.GreaterOrEqual(t, time.Since(start), 100*time.Millisecond)
}

func TestSendEmptyInput(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()
	go l.Accept()

	client := NewTelnetClient(l.Addr().String(), time.Second,
		io.NopCloser(&bytes.Buffer{}),
		&bytes.Buffer{},
	)
	require.NoError(t, client.Connect())
	defer client.Close()
	require.NoError(t, client.Send())
}

func TestReceiveEOF(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)
		conn.Close()
	}()

	buf := &bytes.Buffer{}
	client := NewTelnetClient(l.Addr().String(), time.Second,
		io.NopCloser(&bytes.Buffer{}), buf,
	)
	require.NoError(t, client.Connect())
	defer client.Close()
	require.NoError(t, client.Receive())
	require.Empty(t, buf.String())
}

func TestCloseMany(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()
	go l.Accept()

	client := NewTelnetClient(l.Addr().String(), time.Second,
		io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{},
	)
	require.NoError(t, client.Connect())
	require.NoError(t, client.Close())
	err = client.Close()
	require.Error(t, err)
}
