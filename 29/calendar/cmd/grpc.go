/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	grpcService "github.com/mitrickx/otus-golang-2019/29/calendar/internal/grpc"
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run grpc homework simple entities service",
	Long:  `Run grpc homework simple entities service on port defined in logger`,
	Run: func(cmd *cobra.Command, args []string) {
		runGrpcService()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}

func runGrpcService() {
	// read port
	ports := viper.GetStringMapString("port")

	port, ok := ports["grpc"]
	if !ok {
		port = "50051"
	}

	log := logger.GetLogger()

	storage := NewDbStorage()

	err := grpcService.RunService(port, storage, log)
	if err != nil {
		log.Fatalf("can't run grpc service %s\n", err)
	}
}
