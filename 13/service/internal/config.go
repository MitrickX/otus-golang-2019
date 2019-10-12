package internal

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

var Logger *zap.SugaredLogger

func InitLogger(v *viper.Viper) {

	zapConfig := zap.NewProductionConfig()

	loggerCfg := v.GetStringMapString("logger")
	if level, ok := loggerCfg["level"]; ok {
		// if level unrecognized, just will use NewProductionConfig level, so ignore error
		_ = zapConfig.Level.UnmarshalText([]byte(level))
	}

	var err error
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Can't init zap logger: %s\n", err)
	}

	Logger = logger.Sugar()
}
