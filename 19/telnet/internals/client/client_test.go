package client

import (
	"github.com/mitrickx/otus-golang-2019/19/telnet/internals/test"
	"net"
	"testing"
	"time"
)

// Test client mock
type testClient struct {
	Client
}

// Constructor
func newTestClient(cfg *Config, conn net.Conn) *testClient {
	client := NewClient(cfg)
	client.conn = conn
	return &testClient{
		Client: *client,
	}
}

// Connect (we already has connector from constructor)
func (client *testClient) connect() error {
	return nil
}

// Telnet client send console input to server
func TestInput(t *testing.T) {

	console := test.NewTestConsole()

	defer func() {
		_ = console.Close()
	}()

	clientPipe, serverPipe := net.Pipe()

	// defer closer
	defer func() {
		_ = clientPipe.Close()
		_ = serverPipe.Close()
	}()

	client := newTestClient(&Config{
		Reader: console,
		Writer: console,
	}, clientPipe)

	defer func() {
		client.Stop()
	}()

	go func() {
		err := client.Run()
		if err != nil {
			t.Errorf("telnet client run failed %s\n", err)
			return
		}
	}()

	testMessage := "Hello\nTest telnet client, send console input into server\n123123"
	_, _ = console.Write([]byte(testMessage))

	serverReadBuf := make([]byte, 1024)
	n, err := serverPipe.Read(serverReadBuf)

	if err != nil {
		t.Errorf("Reading from server return error %s\n", err)
	}

	if n == 0 {
		t.Errorf("Reading from server return 0 bytes\n")
	}

	if testMessage != string(serverReadBuf[0:n]) {
		t.Errorf("Reading from server return unexpected result. Expected `%s`\nReturned `%s`\n", testMessage, string(serverReadBuf[0:n]))
	}

}

// Test server send message and message reached console output
func TestOutput(t *testing.T) {
	console := test.NewTestConsole()

	defer func() {
		// clear listener
		console.ListenOnWrite(nil)
		// close console
		_ = console.Close()
	}()

	clientPipe, serverPipe := net.Pipe()

	defer func() {
		_ = clientPipe.Close()
		_ = serverPipe.Close()
	}()

	client := newTestClient(&Config{
		Reader: console,
		Writer: console,
	}, clientPipe)

	defer func() {
		client.Stop()
	}()

	go func() {
		err := client.Run()
		if err != nil {
			t.Errorf("telnet client run failed %s\n", err)
			return
		}
	}()

	written := make(chan struct{}, 1)

	testMessage := "Hello\nTest telnet client, send server input into console\n123123"

	// when message written to console we do assertion
	console.ListenOnWrite(func() {
		consoleReadBuf := make([]byte, 1024)
		n, err := console.Read(consoleReadBuf)

		if err != nil {
			t.Errorf("Reading from console return error %s\n", err)
		}

		if n == 0 {
			t.Errorf("Reading from console return 0 bytes\n")
		}

		if testMessage != string(consoleReadBuf[0:n]) {
			t.Errorf("Reading from console return unexpected result. Expected `%s`\nReturned `%s`\n", testMessage, string(consoleReadBuf[0:n]))
		}

		// signal that is all OK
		written <- struct{}{}
	})

	_, _ = serverPipe.Write([]byte(testMessage))

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Console output not recieved server message after 3 seconds message was sent")
	case <-written:
		// OK
	}

}

// Test input EOF (Ctrl+D), telnet client must stop
func TestEOFInput(t *testing.T) {
	console := test.NewTestConsole()

	defer func() {
		_ = console.Close()
	}()

	clientPipe, serverPipe := net.Pipe()

	// defer closer
	defer func() {
		_ = clientPipe.Close()
		_ = serverPipe.Close()
	}()

	client := newTestClient(&Config{
		Reader: console,
		Writer: console,
	}, clientPipe)

	defer func() {
		client.Stop()
	}()

	stopped := make(chan struct{}, 1)

	go func() {
		err := client.Run()
		stopped <- struct{}{}
		if err != nil {
			t.Errorf("telnet client run failed %s\n", err)
			return
		}
	}()

	testMessage := "Hello\nTest telnet client, send console input into server\n123123"
	_, _ = console.Write([]byte(testMessage))

	serverReadBuf := make([]byte, 1024)
	n, err := serverPipe.Read(serverReadBuf)

	if err != nil {
		t.Errorf("Reading from server return error %s\n", err)
	}

	if n == 0 {
		t.Errorf("Reading from server return 0 bytes\n")
	}

	if testMessage != string(serverReadBuf[0:n]) {
		t.Errorf("Reading from server return unexpected result. Expected `%s`\nReturned `%s`\n", testMessage, string(serverReadBuf[0:n]))
	}

	// Emulating Send EOF
	_ = console.Close()

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Telnet client not stopped after 3 seconds we has sent EOF to console")
	case <-stopped:
		// OK
	}

}

// Test explicit call stop on telnet client
func TestStop(t *testing.T) {
	console := test.NewTestConsole()

	defer func() {
		_ = console.Close()
	}()

	clientPipe, serverPipe := net.Pipe()

	// defer closer
	defer func() {
		_ = clientPipe.Close()
		_ = serverPipe.Close()
	}()

	client := newTestClient(&Config{
		Reader: console,
		Writer: console,
	}, clientPipe)

	defer func() {
		client.Stop()
	}()

	stopped := make(chan struct{}, 1)

	go func() {
		err := client.Run()
		stopped <- struct{}{}
		if err != nil {
			t.Errorf("telnet client run failed %s\n", err)
			return
		}
	}()

	testMessage := "Hello\nTest telnet client, send console input into server\n123123"
	_, _ = console.Write([]byte(testMessage))

	serverReadBuf := make([]byte, 1024)
	n, err := serverPipe.Read(serverReadBuf)

	if err != nil {
		t.Errorf("Reading from server return error %s\n", err)
	}

	if n == 0 {
		t.Errorf("Reading from server return 0 bytes\n")
	}

	if testMessage != string(serverReadBuf[0:n]) {
		t.Errorf("Reading from server return unexpected result. Expected `%s`\nReturned `%s`\n", testMessage, string(serverReadBuf[0:n]))
	}

	// explicit stop
	client.Stop()

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Telnet client not stopped after 3 seconds Stop() method was called")
	case <-stopped:
		// OK
	}

}

// Test case when server close connection
func TestServerCloseConnection(t *testing.T) {
	console := test.NewTestConsole()

	defer func() {
		_ = console.Close()
	}()

	clientPipe, serverPipe := net.Pipe()

	// defer closer
	defer func() {
		_ = clientPipe.Close()
		_ = serverPipe.Close()
	}()

	client := newTestClient(&Config{
		Reader: console,
		Writer: console,
	}, clientPipe)

	defer func() {
		_ = console.Close()
	}()

	stopped := make(chan struct{}, 1)

	go func() {
		err := client.Run()
		stopped <- struct{}{}
		if err != nil {
			t.Errorf("telnet client run failed %s\n", err)
			return
		}
	}()

	testMessage := "Hello\nTest telnet client, send server input into console\n123123"

	// when message written to console we do assertion
	console.ListenOnWrite(func() {
		consoleReadBuf := make([]byte, 1024)
		n, err := console.Read(consoleReadBuf)

		if err != nil {
			t.Errorf("Reading from console return error %s\n", err)
		}

		if n == 0 {
			t.Errorf("Reading from console return 0 bytes\n")
		}

		if testMessage != string(consoleReadBuf[0:n]) {
			t.Errorf("Reading from console return unexpected result. Expected `%s`\nReturned `%s`\n", testMessage, string(consoleReadBuf[0:n]))
		}
	})

	// server send message (write into connection)
	_, _ = serverPipe.Write([]byte(testMessage))

	// server close connection
	_ = serverPipe.Close()

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Telnet client not stopped after 3 seconds server close connection")
	case <-stopped:
		// OK
	}
}
