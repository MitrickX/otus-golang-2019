package tests

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
)

var runnerOptions = godog.Options{
	Format:    "pretty", // progress, pretty
	Paths:     []string{"../features/"},
	Randomize: 0,
}

// Test entry point
func TestMain(m *testing.M) {

	features := flag.String("features", "", `-features="create_event,delete_event"`)
	flag.Parse()

	if *features != "" {
		featureList := strings.Split(*features, ",")
		pathPrefix := "../features/"
		var paths []string
		for _, f := range featureList {
			paths = append(paths, pathPrefix+f+".feature")
		}
		runnerOptions.Paths = paths
	}

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		t := newFeatureTest()
		FeatureContext(s, t)
	}, runnerOptions)

	os.Exit(status)
}
