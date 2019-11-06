package client

import (
	"context"
	"github.com/mitrickx/otus-golang-2019/19/telnet/internals/copy"
	"io"
	"log"
	"net"
)

// Simple telnet client
type Client struct {
	cfg    *Config
	conn   net.Conn
	cancel context.CancelFunc // for cancel copy processing
}

// Simple telnet client constructor
func NewClient(cfg *Config) *Client {
	client := &Client{
		cfg: cfg,
	}
	return client
}

// dial (connect with remote)
// set up struct conn
func (client *Client) connect() error {

	if client.conn != nil {
		return nil
	}

	if client.cfg.Logging {
		log.Printf("telnet client is connecting to %s...\n", client.cfg.Address)
	}

	dialer := net.Dialer{
		Timeout: client.cfg.Timeout,
	}

	conn, err := dialer.Dial(client.cfg.Network, client.cfg.Address)
	if err != nil {
		return err
	}

	client.conn = conn

	return nil
}

// disconnect if was connected
func (client *Client) disconnect() {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil && client.cfg.Logging {
			log.Println(err)
		}
	}
}

// Run simple telnet client
// Connect by tcp/tcp4/tcp6 to Address
// Read from and write to connection
// io.EOF stops running
func (client *Client) Run() error {

	err := client.connect()
	if err != nil {
		return err
	}

	// init our context with cancel
	ctx, cancelFunc := context.WithCancel(context.Background())
	client.cancel = cancelFunc

	// read from connection to stdout (or any Writer)
	readResultCh := make(chan error)
	copy.RunCopier(ctx, client.conn, client.cfg.Writer, readResultCh)

	// write to connection from stdin (or any Reader)
	writeResultCh := make(chan error)
	copy.RunCopier(ctx, client.cfg.Reader, client.conn, writeResultCh)

	var resultError error

	// select first error and Stop our copy go-routines (with context cancel function)
	select {
	case <-ctx.Done():
		client.disconnect()
	case resultError = <-readResultCh:
		client.disconnect()
		if resultError == io.EOF && client.cfg.Logging {
			log.Println("remote Address close connection")
		}
	case resultError = <-writeResultCh:
		client.disconnect()
		if resultError == io.EOF && client.cfg.Logging {
			log.Println("close connection")
		}
	}

	// no need to return EOF, EOF for inner logic
	if resultError == io.EOF {
		resultError = nil
	}

	return resultError

}

// Stop simple telnet client
func (client *Client) Stop() {
	if client.cancel != nil {
		if client.cfg.Logging {
			log.Println("stopping telnet client")
		}
		client.cancel()
	}
}
