package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	exitCodePtr := flag.Int("code", 0, "exit code of program")

	flag.Parse()

	// show all env vars
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}

	os.Exit(*exitCodePtr)
}
