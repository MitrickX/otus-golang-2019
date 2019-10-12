package cmd

import (
	"github.com/mitrickx/otus-golang-2019/13/service/internal"
	"github.com/spf13/viper"
	"net/http"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run empty http service on port declared in config",
	Long:  `Run empty http service on port declared in config, see ./configs/config.yaml. Default port is 8080`,
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
		internal.Logger.Info("Request processed")
	})

	// try run service and failed log fatal
	internal.Logger.Info("Run server on port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		internal.Logger.Fatal("Run server error: ", err)
	}
}
