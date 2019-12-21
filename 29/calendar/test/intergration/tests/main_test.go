package tests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("Wait for 5s for services...")
	time.Sleep(5 * time.Second)

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty", // progress, pretty
		Paths:     []string{"../features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
