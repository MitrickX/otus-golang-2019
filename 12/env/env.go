package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// show all env vars
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}
}
