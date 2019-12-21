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
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/logger"
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/notificaiton"
	"github.com/spf13/cobra"
)

// senderCmd represents the sender command
var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "A notification sender",
	Long:  `A notification sender (just print enqueued events in log)`,
	Run: func(cmd *cobra.Command, args []string) {
		runNotificationSender()
	},
}

func init() {
	rootCmd.AddCommand(senderCmd)
}

// Run notification sender (this sender just print into log)
func runNotificationSender() {
	log := logger.GetLogger()
	queue := NewNotificationQueue()
	sender := notificaiton.NewLogSender(queue, *log)
	err := sender.Run()
	if err != nil {
		log.Fatal("can't run sender, error happened: %s", err)
	}
}
