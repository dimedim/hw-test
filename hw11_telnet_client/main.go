package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "server conn timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "not enough arguments. Usage: %s [--timeout duration] host port\n", os.Args[0])
		return
	}
	in := os.Stdin
	out := os.Stdout
	address := net.JoinHostPort(args[0], args[1])

	telCl := NewTelnetClient(address, *timeout, in, out)

	fmt.Fprintf(os.Stderr, "...Connect to %s with timeout %v\n", address, *timeout)
	if err := telCl.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "...Connection error: %v\n", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	doneSend := make(chan struct{})
	doneReceive := make(chan struct{})
	go func() {
		if err := telCl.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "...send error: %s\n", err.Error())
		}
		close(doneSend)
	}()
	go func() {
		if err := telCl.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "...receive error: %s\n", err.Error())
		}
		close(doneReceive)
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "...Interrupted, exiting")
	case <-doneSend:
		fmt.Fprintln(os.Stderr, "...EOF")
	case <-doneReceive:
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	}

	err := telCl.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "...close telenet client error: %s\n", err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, "...closed, OK")
}
