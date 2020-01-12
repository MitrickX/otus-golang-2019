package cmd

import (
	"log"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/logger"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/monitoring"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/notificaiton"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/memory"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/sql"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	cfgDefaultFilePath            = "./configs/config.yaml"
	defaultSqlMetricsExporterPort = "9103"
)

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")
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
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
}

func initLogger() {
	logger.InitLogger(viper.GetViper())
}

func NewDbStorage() entities.Storage {
	log := logger.GetLogger()

	var storage entities.Storage

	dbConf := viper.GetStringMapString("db")

	if dbConf == nil {
		storage = memory.NewStorage()
	} else {
		dbConfig, err := sql.NewConfig(dbConf)
		if err != nil {
			log.Fatalf("can't init sql storage %s\n", err)
		}
		storage, err = sql.NewStorage(*dbConfig)
		if err != nil {
			log.Fatalf("can't init sql storage %s\n", err)
		}
	}

	return storage
}

func NewSqlMetrics(storage *sql.Storage) (*monitoring.SqlMetrics, error) {

	log := logger.GetLogger()

	dbConf := viper.GetStringMap("db")

	exporterPort := defaultSqlMetricsExporterPort

	prometheusConfigValue, ok := dbConf["prometheus"]
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

	return monitoring.NewSqlMetrics(storage, exporterPort, log)
}

func NewNotificationQueue() notificaiton.Queue {
	log := logger.GetLogger()

	nConf := viper.GetStringMap("notification")
	if nConf == nil {
		log.Fatal("can't init queue, notification settings not found in `notification` key of config")
	}

	confVar, ok := nConf["queue"]
	if !ok {
		log.Fatal("can't init queue, queue settings not found in `queue` key key of `notification` config")
	}

	qConf := cast.ToStringMapString(confVar)

	qConfig, err := notificaiton.NewConfig(qConf)
	if err != nil {
		log.Fatalf("can't init queue %s\n", err)
	}

	queue, err := notificaiton.NewRabbitQueue(*qConfig)
	if err != nil {
		log.Fatalf("can't init queue %s\n", err)
	}

	return queue
}
