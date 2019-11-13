package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

var logger *zap.SugaredLogger

func InitLogger(v *viper.Viper) {

	zapConfig := zap.NewProductionConfig()

	loggerCfg := v.GetStringMapString("logger")
	if level, ok := loggerCfg["level"]; ok {
		// if level unrecognized, just will use NewProductionConfig level, so ignore error
		_ = zapConfig.Level.UnmarshalText([]byte(level))
	}

	var err error
	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Can't init zap logger: %s\n", err)
	}

	logger = zapLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		log.Fatal("Logger is not inited")
	}
	return logger
}
