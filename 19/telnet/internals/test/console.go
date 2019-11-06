package test

import (
	"io"
	"sync"
)

// Simple structure to mock terminal
// Basically it just reader/writer/closer
type testConsole struct {
	buf      []byte
	bufMx    *sync.RWMutex
	closed   bool
	closedMx *sync.RWMutex
	onWrite  func()
}

// Constructor
func NewTestConsole() *testConsole {
	return &testConsole{
		bufMx:    &sync.RWMutex{},
		closedMx: &sync.RWMutex{},
	}
}

// Goroutine safe isClosed checker
func (console *testConsole) isClosed() bool {
	console.closedMx.RLock()
	closed := console.closed
	console.closedMx.RUnlock()
	return closed
}

// Read from test console
func (console *testConsole) Read(p []byte) (int, error) {

	if console.isClosed() {
		return 0, io.EOF
	}

	console.bufMx.RLock()
	n := copy(p, console.buf)
	console.bufMx.RUnlock()

	return n, nil
}

// Write to test console
func (console *testConsole) Write(p []byte) (int, error) {

	if console.isClosed() {
		return 0, io.EOF
	}

	console.bufMx.Lock()
	console.buf = append(console.buf, p...)
	console.bufMx.Unlock()

	if console.onWrite != nil {
		console.onWrite()
	}

	return len(p), nil
}

// Set listener on onWrite event
func (console *testConsole) ListenOnWrite(fn func()) {
	console.onWrite = fn
}

// Clear all console (clear inner buffer)
func (console *testConsole) Clear() {
	console.bufMx.Lock()
	console.buf = console.buf[:0]
	console.bufMx.Unlock()
}

// Close console (emulate io.EOF situation)
func (console *testConsole) Close() error {
	console.Clear()
	console.closedMx.Lock()
	console.closed = true
	console.closedMx.Unlock()
	return nil
}
