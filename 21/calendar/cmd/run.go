package cmd

import (
	httpService "github.com/mitrickx/otus-golang-2019/21/calendar/internal/http"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run http homework simple calendar service",
	Long:  `Run http homework simple calendar service on port defined in logger`,
	Run: func(cmd *cobra.Command, args []string) {
		runHttpService()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

}

func runHttpService() {

	// read port
	port := viper.GetString("port")
	if port == "" {
		port = "8080"
	}

	serviceLogger := logger.GetLogger()
	httpService.RunService(port, serviceLogger)
}
