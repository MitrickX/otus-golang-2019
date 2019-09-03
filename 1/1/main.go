package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"math/rand"
	"time"
)

// List of possible ntp hosts
var hosts = []string{
	"ntp1.stratum1.ru",
	"ntp2.stratum1.ru",
	"ntp3.stratum1.ru",
	"ntp4.stratum1.ru",
	"ntp5.stratum1.ru",
}

// Max attempts for define current time by ntp host(s)
// See getCurrentTime
const attempts = 3

func init() {
	// Init random
	rand.Seed(time.Now().Unix())
}

// Choose randomly host from hosts list
func getHost() string {
	n := len(hosts)
	i := rand.Intn(n)
	return hosts[i]
}

// Get current time from random choosen ntp host
// Will make several attempts (max limit is const attempts) until got first successful one
// After each fail attempt try choose another ntp host
// In most cases first attemp is successful
func getCurrentTime() (time.Time, error) {
	var time time.Time
	var error error
	for attept := 0; attept < attempts; attept++ {
		host := getHost()
		time, error = ntp.Time(host)
		if error != nil {
			break
		}
	}
	return time, error
}

//
// Define current time and print it
// On fail print about reason of fail
//
func main() {
	now, err := getCurrentTime()
	if err != nil {
		fmt.Printf("Couldn't define current datetime.\nReason is %s\n", err)
	} else {
		fmt.Printf("Current datetime is %s\n", now.Format(time.RFC1123Z))
	}
}
