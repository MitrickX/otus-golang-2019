package tests

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

const cfgFilePath = "../../../../configs/config.yaml"

type IntegrationTestsConfig struct {
	DbStorage        *sql.Storage
	SenderOutputPath string
	RunnerPaths      []string
}

var testConfig *IntegrationTestsConfig

func init() {

	cfgPath := flag.String("config", "", `--config=<path>`)

	features := flag.String("features", "", `-features="create_event,delete_event"`)
	featuresPath := flag.String("features-path", "", `-features-path="./features/"`)

	flag.Parse()

	if *cfgPath == "" {
		*cfgPath = cfgFilePath
	}

	viper.SetConfigFile(*cfgPath)

	// If a logger file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}

	testConfig = new(IntegrationTestsConfig)

	if *features != "" {
		featureList := strings.Split(*features, ",")
		pathPrefix := "../features/"
		if *featuresPath != "" {
			pathPrefix = *featuresPath
		}
		var paths []string
		for _, f := range featureList {
			paths = append(paths, pathPrefix+f+".feature")
		}
		testConfig.RunnerPaths = paths

		log.Println("run on features", paths)
	} else if *featuresPath != "" {
		testConfig.RunnerPaths = append(testConfig.RunnerPaths, *featuresPath)
		log.Println("run on features", testConfig.RunnerPaths)
	}

	outputPath := os.Getenv("SENDER_OUTPUT_PATH")
	if outputPath != "" {
		testConfig.SenderOutputPath = outputPath
		log.Printf("Sender output path set to %s\n", outputPath)
	}

	dbStorage, err := newDbStorage(viper.GetViper())
	if err != nil {
		log.Fatal(err)
	}

	testConfig.DbStorage = dbStorage

}

func GetTestConfig() *IntegrationTestsConfig {
	if testConfig == nil {
		log.Fatal("Integration test config is not initialized!")
	}
	return testConfig
}

func newDbStorage(v *viper.Viper) (*sql.Storage, error) {

	dbConf := v.GetStringMapString("db")

	cfg, err := sql.NewConfig(dbConf)
	if err != nil {
		return nil, err
	}

	log.Println("wait for connect to db...")

	storage, err := sql.NewStorage(*cfg)
	if err != nil {
		return nil, err
	}

	log.Println("successfully connected to db")

	return storage, nil

}
