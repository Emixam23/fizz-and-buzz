package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/utils/logger"
)

// Config is the configuration container of the app
type Config struct {
	Logger *logger.Config

	DatabaseConfig   *infra.Config
	FnbServiceConfig *fnbservice.Config
	RestAPIConfig    *ui.Config
}

const (
	envLoggerLevel        = "LOGGER_LEVEL"
	envLoggerAsJSONFormat = "LOGGER_AS_JSON_FORMAT"

	envDatabaseHost              = "DATABASE_HOST"
	envDatabaseName              = "DATABASE_NAME"
	envDatabaseUser              = "DATABASE_USER"
	envDatabasePassword          = "DATABASE_PASSWORD"
	envDatabasePort              = "DATABASE_PORT"
	envDatabaseTimeoutSeconds    = "DATABASE_TIMEOUT_SECONDS"
	envDatabaseRetryAmountOnFail = "DATABASE_RETRY_AMOUNT_ON_FAIL"

	envFnbServiceZero = "FNB_SERVICE_ZERO"

	envRouterHost = "ROUTER_HOST"
	envRouterPort = "ROUTER_PORT"
	envRouterMode = "ROUTER_MODE"
)

const defaultLoggerLevel = zerolog.InfoLevel

// LoadConfig loads the config based on the provided env keys/values
func LoadConfig() *Config {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig() // As this is just optional, it's ok  if an error happens  because no file were retrieved

	viper.AutomaticEnv()

	viper.SetDefault(envLoggerLevel, defaultLoggerLevel)
	viper.SetDefault(envLoggerAsJSONFormat, true)
	viper.SetDefault(envDatabaseTimeoutSeconds, 2)
	viper.SetDefault(envFnbServiceZero, 1)
	viper.SetDefault(envRouterHost, "0.0.0.0")
	viper.SetDefault(envRouterMode, gin.ReleaseMode)
	viper.SetDefault(envRouterPort, "8080")

	loggerLevel, err := zerolog.ParseLevel(viper.GetString(envLoggerLevel))
	if err != nil {
		loggerLevel = defaultLoggerLevel
	}
	routerMode := viper.GetString(envRouterMode)

	appConfig := Config{
		Logger: &logger.Config{
			JSON:  viper.GetBool(envLoggerAsJSONFormat),
			Level: loggerLevel,
		},
		DatabaseConfig: &infra.Config{
			Host:                  viper.GetString(envDatabaseHost),
			DatabaseName:          viper.GetString(envDatabaseName),
			User:                  viper.GetString(envDatabaseUser),
			Password:              viper.GetString(envDatabasePassword),
			Port:                  viper.GetUint64(envDatabasePort),
			DefaultTimeoutSeconds: viper.GetUint(envDatabaseTimeoutSeconds),
			RetryAmountOnFail:     viper.GetUint32(envDatabaseRetryAmountOnFail),
		},
		FnbServiceConfig: &fnbservice.Config{
			Zero: viper.GetUint64(envFnbServiceZero),
		},
		RestAPIConfig: &ui.Config{
			Host: viper.GetString(envRouterHost),
			Mode: routerMode,
			Port: viper.GetUint32(envRouterPort),
		},
	}

	return &appConfig

}
