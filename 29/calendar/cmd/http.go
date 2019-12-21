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
	httpService "github.com/mitrickx/otus-golang-2019/29/calendar/internal/http"
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/logger"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run http homework simple entities service.",
	Long:  `Run http homework simple entities service on port defined in logger.`,
	Run: func(cmd *cobra.Command, args []string) {
		runHttpService()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func runHttpService() {
	// read port

	ports := viper.GetStringMapString("port")

	port, ok := ports["http"]
	if !ok {
		port = "8080"
	}

	log := logger.GetLogger()

	storage := NewDbStorage()

	err := httpService.RunService(port, storage, log)
	if err != nil {
		log.Fatalf("can't run http service %s\n", err)
	}
}
