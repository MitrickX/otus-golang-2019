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
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/logger"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/notificaiton"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"time"

	"github.com/spf13/cobra"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "A notification scheduler",
	Long:  `A notification scheduler.`,
	Run: func(cmd *cobra.Command, args []string) {
		runNotificationScheduler()
	},
}

func init() {
	rootCmd.AddCommand(schedulerCmd)
}

// Run notification scheduler
func runNotificationScheduler() {

	log := logger.GetLogger()

	nConf := viper.GetStringMap("notification")
	if nConf == nil {
		log.Fatal("can't init scheduler, notification settings not found in `notification` key of config")
	}

	confVar, ok := nConf["scheduler"]
	if !ok {
		log.Fatal("can't init scheduler, queue settings not found in `scheduler` key key of `notification` config")
	}

	sConf := cast.ToStringMapString(confVar)

	scanTimeoutVal, ok := sConf["scan_timeout"]
	if !ok {
		scanTimeoutVal = "1m"
	}

	scanTimeout, err := time.ParseDuration(scanTimeoutVal)
	if err != nil {
		log.Fatal("can't init scheduler, fail on parsing `scan_timeout` value == `%s`", scanTimeoutVal)
	}

	queue := NewNotificationQueue()
	storage := NewDbStorage()

	scheduler := notificaiton.NewScheduler(
		scanTimeout,
		storage,
		queue,
		log,
	)

	err = scheduler.Run()
	if err != nil {
		log.Fatal("can't run scheduler, error happened: %s", err)
	}

}
