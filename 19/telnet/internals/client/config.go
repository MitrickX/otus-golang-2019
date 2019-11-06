package client

import (
	"io"
	"time"
)

// Config for simple telnet client
type Config struct {
	Network string
	Address string
	Timeout time.Duration
	Reader  io.Reader // by default here is os.Stdin
	Writer  io.Writer // by default here is os.Stdout
	Logging bool      // need Logging
}
