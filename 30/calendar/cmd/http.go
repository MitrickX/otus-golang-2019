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
	httpService "github.com/mitrickx/otus-golang-2019/30/calendar/internal/http"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/logger"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	defaultHttpPort                = "8080"
	defaultHttpMetricsExporterPort = "9102"
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
	// read port from config

	port := defaultHttpPort

	httpConfig := viper.GetStringMap("http")
	portValue, ok := httpConfig["port"]
	if ok {
		portVal, ok := portValue.(string)
		if ok {
			port = portVal
		}
	}

	exporterPort := defaultHttpMetricsExporterPort

	prometheusConfigValue, ok := httpConfig["prometheus"]
	prometheusConfig := cast.ToStringMap(prometheusConfigValue)
	if ok {
		portValue, ok := prometheusConfig["port"]
		if ok {
			portVal, ok := portValue.(string)
			if ok {
				exporterPort = portVal
			}
		}
	}

	log := logger.GetLogger()

	storage := NewDbStorage()

	err := httpService.RunService(port, storage, log, exporterPort)
	if err != nil {
		log.Fatalf("can't run http service %s\n", err)
	}
}
