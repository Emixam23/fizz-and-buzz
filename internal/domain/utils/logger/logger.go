package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// Config contains the configuration for the logger
type Config struct {
	JSON  bool
	Level zerolog.Level
}

// SeverityHook duplicates the level attribute into a severity attribute.
// This is done because of GCP Monitoring tool that only understands the severity attribute.
type SeverityHook struct{}

// Run is a simple hook that will add a "severity" key and the related log level as value to it when Log...Msg() is triggered
func (SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", level.String())
	}
}

// New creates and initializes a new zerolog logger
func New(config *Config) (logger zerolog.Logger) {

	if config == nil {
		config = &Config{
			JSON:  true,
			Level: zerolog.InfoLevel,
		}
	}

	if config.JSON {
		logger = zerolog.New(os.Stdout)
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return logger.Hook(SeverityHook{}).
		With().
		Caller().
		Timestamp().
		Logger().
		Level(config.Level)
}
