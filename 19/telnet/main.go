package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config for simple telnet client
type config struct {
	network string
	address string
	timeout time.Duration
}

// Parse cli arguments to config
func parseArgs() config {
	cfg := config{}
	var timeout string

	flag.StringVar(&cfg.network, "network", "tcp", "tcp, tcp4, tcp6, default is tcp")
	flag.StringVar(&cfg.address, "address", "", "string, address in format host:port, required")
	flag.StringVar(&timeout, "timeout", "10s", "string, timeout for connect to server, default is 10s")
	flag.Parse()

	if cfg.address == "" {
		flag.Usage()
		os.Exit(1)
	}

	// parse string timeout to duration
	d, err := time.ParseDuration(timeout)
	if err != nil {
		d = 10 * time.Second
	}
	cfg.timeout = d

	return cfg
}

// Copy, similar with io.copyBuffer except it actually break loop on io.EOF (in this telnet client it is significant)
func ioCopy(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {

	for {
		nr, er := src.Read(buf)

		// break on any read error (and also on EOF)
		if er != nil {
			err = er
			break
		}

		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

	}
	return written, err
}

// Run ioCopy go-routine
// Params
//  - context for canceling
//  - reader from where will read
//  - writer where will write
//  - channel of error that could be happened
func runCopy(ctx context.Context, reader io.Reader, writer io.Writer, result chan<- error) {
	go func() {

		// buffer for ioCopy
		buf := make([]byte, 4096)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				result <- nil
				break LOOP
			default:
				// writer <- reader
				_, err := ioCopy(writer, reader, buf)
				if err != nil {
					result <- err
					break LOOP
				}
			}
		}
	}()
}

// Run simple telnet client
// Connect by tcp/tcp4/tcp6 to address (passed as config argument)
// Read from and write to connection
// Ctrl+D stops telnet client (works for Linux)
// Ctrl+C stops telnet client (works for Windows and Linux)
// Params
// - config
// - os.Signal channel for listening interruption/termination signal and close connection gracefully
func telnet(cfg config, sigCh chan os.Signal) (resultError error) {

	dialer := net.Dialer{
		Timeout: cfg.timeout,
	}

	log.Printf("Connecting to %s...\n", cfg.address)
	connection, err := dialer.Dial(cfg.network, cfg.address)

	if err != nil {
		resultError = err
		return
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			resultError = err
		}
	}()

	log.Println("Connected")

	// init our context with cancel
	ctx, cancelFunc := context.WithCancel(context.Background())

	// read from connection to stdout
	readResultCh := make(chan error)
	runCopy(ctx, connection, os.Stdout, readResultCh)

	// write to connection from stdin
	writeResultCh := make(chan error)
	runCopy(ctx, os.Stdin, connection, readResultCh)

	// select first error and stop our ioCopy go-routines (with context cancel function)
	select {
	case err := <-readResultCh:
		cancelFunc()
		if err == io.EOF {
			log.Println("Remote address close connection")
		} else if err != nil {
			resultError = err
		}
	case err := <-writeResultCh:
		cancelFunc()
		if err == io.EOF {
			log.Println("Close connection")
		} else if err != nil {
			resultError = err
		}
	case <-sigCh:
		cancelFunc()
	}

	return

}

// Entry point, parse args to config and call run telnet client
func main() {
	cfg := parseArgs()

	// listen interruption/termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	err := telnet(cfg, sigCh)
	if err != nil {
		log.Fatal(err)
	}

}
