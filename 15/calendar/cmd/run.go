package cmd

import (
	"github.com/mitrickx/otus-golang-2019/15/calendar/internals/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
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

	// main service handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.GetLogger().Info("Request processed")
	})

	// try run service and failed log fatal
	logger.GetLogger().Info("Run server on port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.GetLogger().Fatal("Run server error: ", err)
	}
}
