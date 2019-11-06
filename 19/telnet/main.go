package main

import (
	"flag"
	"github.com/mitrickx/otus-golang-2019/19/telnet/internals/client"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Parse some cli arguments to Config options
func parseArgs(cfg *client.Config) {
	var timeout string

	flag.StringVar(&cfg.Network, "network", "tcp", "tcp, tcp4, tcp6, default is tcp")
	flag.StringVar(&cfg.Address, "address", "", "string, address in format host:port, required")
	flag.StringVar(&timeout, "timeout", "10s", "string, timeout for connect to server, default is 10s")
	flag.Parse()

	if cfg.Address == "" {
		flag.Usage()
		os.Exit(1)
	}

	// parse string timeout to duration
	d, err := time.ParseDuration(timeout)
	if err != nil {
		d = 10 * time.Second
	}
	cfg.Timeout = d

}

// Entry point, parse args to Config and call Run telnet Client
// Ctrl+C stops telnet Client (works for Windows and Linux)
// Ctrl+D stops telnet Client (works for Linux)
func main() {
	//
	cfg := &client.Config{
		Logging: true,
	}

	// parse arguments and write options into Config
	parseArgs(cfg)

	// set up reader/write (by default is stdin/stdout)
	cfg.Reader = os.Stdin
	cfg.Writer = os.Stdout

	telnet := client.NewClient(cfg)

	sigCh := make(chan os.Signal, 1) // for graceful shutdown on syscall.SIGINT, syscall.SIGTERM or just close telnet by method
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		telnet.Stop()
	}()

	err := telnet.Run()
	if err != nil {
		log.Fatal(err)
	}

}
