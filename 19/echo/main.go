package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// Config for echo server
type config struct {
	address string
	timeout int // in ms
}

// new line (default is \n, for windows is \r\n)
var newLine = []byte{'\n'}

// init of package
func init() {
	if runtime.GOOS == "windows" {
		newLine = []byte{'\r', '\n'}
	}
}

// Parse cli arguments to config
func parseArgs() config {
	cfg := config{}
	flag.StringVar(&cfg.address, "address", "127.0.0.1:9000", "string, address, default is 127.0.0.1:8080")
	flag.IntVar(&cfg.timeout, "timeout", 0, "int, timeout of ticker, default is 0 ms (without ticker)")
	flag.Parse()
	return cfg
}

// Ticker worker for echo server
// Params
// - context for stop worker
// - connection
// - timeout
// - error channel of goroutine of worker
func ticker(ctx context.Context, connection net.Conn, timeout int, result chan<- error) {
	go func() {
		ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	LOOP:
		for {
			select {
			case <-ctx.Done():
				result <- nil
				break LOOP
			case t := <-ticker.C:
				message := t.Format("15:04:05.000")
				_, err := connection.Write(append([]byte(message), newLine...))
				if err != nil {
					result <- err
					break LOOP
				}
			}
		}
	}()
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

// Echo worker for echo server
// Params
// - context for stop worker
// - connection
// - error channel of goroutine of worker
func echo(ctx context.Context, connection net.Conn, result chan<- error) {
	go func() {

		//reader := newConnReader(connection)
		//var output []byte
		buffer := make([]byte, 4096)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			default:
				_, err := ioCopy(connection, connection, buffer)
				if err != nil {
					result <- err
					break LOOP
				}

			}
		}
	}()
}

//
// Handle connection
// Run echo worker
// Run (or not) ticker
//
// Params
// - context for stop handle
// - config to run (or not) ticker
// - connection
func handleConn(handleCtx context.Context, cfg config, connection net.Conn) error {

	ctx, cancel := context.WithCancel(context.Background())

	tickerResultCh := make(chan error)
	echoResultCh := make(chan error)

	if cfg.timeout > 0 {
		ticker(ctx, connection, cfg.timeout, tickerResultCh)
	}

	echo(ctx, connection, echoResultCh)

	select {
	case <-handleCtx.Done():
		cancel()
	case err := <-tickerResultCh:
		cancel()
		if err != nil {
			return err
		}
	case err := <-echoResultCh:
		cancel()
		if err != nil {
			return err
		}
	}

	return nil
}

func closeListener(listener net.Listener) error {
	err := listener.Close()
	log.Println("TCP listener closed")
	return err
}

// Run echo sever listener
func runListener(cfg config, sigCh chan os.Signal) (resultError error) {

	// create listener
	ln, err := net.Listen("tcp", cfg.address)
	if err != nil {
		resultError = err
		return
	}

	// for cancel connection handlers
	ctx, cancelFunc := context.WithCancel(context.Background())

	// interruption/termination signal handler
	go func() {
		<-sigCh      // read interruption/termination signal
		cancelFunc() // for cancel connection handlers and accept loop

		// we must close listener here (not in defer), case Accept is blocking
		resultError = closeListener(ln)
	}()

	// on close lis we must wait for all
	wg := &sync.WaitGroup{}

	// Accept connection loop

LOOP:
	for {

		select {
		case <-ctx.Done():
			break LOOP
		default:

			connection, err := ln.Accept()

			// accept error
			if err != nil {
				log.Println(err)
				cancelFunc()  // cancel connection handlers and accept loop
				continue LOOP // on next loop iteration we done
			}

			// accepted, inc wg for sync closing before return
			log.Println("Connection accepted")
			wg.Add(1)

			// handle connection in goroutine
			go func() {

				err := handleConn(ctx, cfg, connection)
				if err != nil {
					log.Println(err)
				}

				// close connection after handler is done
				err = connection.Close()
				if err != nil {
					log.Println("Connection closed with error", err)
				} else {
					log.Println("Connection closed")
				}

				wg.Done()
			}()
		}
	}

	wg.Wait()

	return
}

// Entry point, parse args to config and call run echo server
func main() {
	cfg := parseArgs()

	// listen interruption/termination signals
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Run echo server on address %s\n", cfg.address)
	err := runListener(cfg, sigCh)
	if err != nil {
		log.Fatal(err)
	}
}
