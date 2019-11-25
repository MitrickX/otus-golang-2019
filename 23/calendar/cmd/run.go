package cmd

import (
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/domain/entities"
	grpcService "github.com/mitrickx/otus-golang-2019/23/calendar/internal/grpc"
	httpService "github.com/mitrickx/otus-golang-2019/23/calendar/internal/http"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/logger"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/storage/memory"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/storage/sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [type=http|grpc]",
	Short: "Run http/grpc homework simple entities service",
	Long:  `Run http/grpc homework simple entities service on port defined in logger`,
	Run: func(cmd *cobra.Command, args []string) {

		log := logger.GetLogger()

		serviceType := ""
		if len(args) <= 0 {
			serviceType = "http"
		} else if args[0] == "http" {
			serviceType = "http"
		} else if args[0] == "grpc" {
			serviceType = "grpc"
		} else {
			serviceType = args[0]
		}

		if serviceType == "http" {
			runHttpService()
		} else if serviceType == "grpc" {
			runGrpcService()
		} else {
			log.Fatalf("unknown service %s\n", args[0])
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runHttpService() {
	// read port
	ports := viper.GetStringMapString("port")

	port, ok := ports["http"]
	if !ok {
		port = "8080"
	}

	log := logger.GetLogger()

	storage := newStorage()

	err := httpService.RunService(port, storage, log)
	if err != nil {
		log.Fatalf("can't run http service %s\n", err)
	}
}

func runGrpcService() {
	// read port
	ports := viper.GetStringMapString("port")

	port, ok := ports["grpc"]
	if !ok {
		port = "50051"
	}

	log := logger.GetLogger()

	storage := newStorage()

	err := grpcService.RunService(port, storage, log)
	if err != nil {
		log.Fatalf("can't run grpc service %s\n", err)
	}
}

func newStorage() entities.Storage {
	log := logger.GetLogger()

	var storage entities.Storage

	dbConf := viper.GetStringMapString("db")
	if dbConf == nil {
		storage = memory.NewStorage()
	} else {
		dbConfig, err := sql.NewConfig(dbConf)
		if err != nil {
			log.Fatalf("can't init sql storage %s\n", err)
		}
		storage, err = sql.NewStorage(*dbConfig)
		if err != nil {
			log.Fatalf("can't init sql storage %s\n", err)
		}
	}

	return storage
}
