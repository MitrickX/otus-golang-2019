package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"

	grpcService "github.com/mitrickx/otus-golang-2019/22/calendar/internal/grpc"
	httpService "github.com/mitrickx/otus-golang-2019/22/calendar/internal/http"
	"github.com/mitrickx/otus-golang-2019/22/calendar/internal/logger"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [type=http|grpc]",
	Short: "Run http/grpc homework simple calendar service",
	Long:  `Run http/grpc homework simple calendar service on port defined in logger`,
	Run: func(cmd *cobra.Command, args []string) {

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
			log.Fatalf("Unknown service %s\n", args[0])
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

	serviceLogger := logger.GetLogger()

	httpService.RunService(port, serviceLogger)
}

func runGrpcService() {
	// read port
	ports := viper.GetStringMapString("port")

	port, ok := ports["grpc"]
	if !ok {
		port = "50051"
	}

	serviceLogger := logger.GetLogger()

	grpcService.RunService(port, serviceLogger)
}
