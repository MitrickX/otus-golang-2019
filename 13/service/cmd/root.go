package cmd

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/13/service/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const cfgDefaultFilePath = "./configs/config.yaml"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "service",
	Short: "A simple homework micro service",
	Long:  `A simple homework micro service, initialize config, logger and run empty http service on port declared in config`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	cfgFilePath := cfgFile
	if cfgFilePath == "" {
		cfgFilePath = cfgDefaultFilePath
	}

	viper.SetConfigFile(cfgFilePath)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLogger() {
	internal.InitLogger(viper.GetViper())
}
