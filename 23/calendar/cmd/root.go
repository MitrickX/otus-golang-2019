package cmd

import (
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string

const cfgDefaultFilePath = "./configs/config.yaml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "entities",
	Short: "Homework simple entities service application",
	Long:  `Homework simple entities service application: create, edit, delete events`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "logger", "", "logger file (default is $HOME/.entities.yaml)")
}

// initConfig reads in logger file and ENV variables if set.
func initConfig() {

	cfgFilePath := cfgFile
	if cfgFilePath == "" {
		cfgFilePath = cfgDefaultFilePath
	}

	viper.SetConfigFile(cfgFilePath)

	// If a logger file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using logger file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
}

func initLogger() {
	logger.InitLogger(viper.GetViper())
}
