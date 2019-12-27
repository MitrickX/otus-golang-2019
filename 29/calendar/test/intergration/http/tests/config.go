package tests

import (
	"log"
	"os"

	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

const cfgFilePath = "../../../../configs/config.yaml"

type IntegrationTestsConfig struct {
	DbStorage        *sql.Storage
	SenderOutputPath string
}

var testConfig *IntegrationTestsConfig

func init() {

	viper.SetConfigFile(cfgFilePath)

	// If a logger file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}

	dbStorage, err := newDbStorage(viper.GetViper())
	if err != nil {
		log.Fatal(err)
	}
	testConfig = &IntegrationTestsConfig{
		DbStorage: dbStorage,
	}
	outputPath := os.Getenv("SENDER_OUTPUT_PATH")
	if outputPath != "" {
		testConfig.SenderOutputPath = outputPath
		log.Printf("Sender output path set to %s\n", outputPath)
	}

}

func GetTestConfig() *IntegrationTestsConfig {
	if testConfig == nil {
		log.Fatal("Integration test config is not initialized!")
	}
	return testConfig
}

func newDbStorage(v *viper.Viper) (*sql.Storage, error) {

	dbConf := v.GetStringMapString("db")

	dbConf["host"] = "localhost"
	dbConf["port"] = "5555"

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
