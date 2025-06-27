package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type params struct {
		config *Config
	}
	tests := []struct {
		name       string
		params     params
		wantLogger zerolog.Logger
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				config: &Config{
					JSON:  true,
					Level: zerolog.DebugLevel,
				},
			},
			wantLogger: zerolog.New(os.Stdout).Hook(SeverityHook{}).With().Caller().Timestamp().Logger().Level(zerolog.DebugLevel),
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				config: &Config{
					JSON:  false,
					Level: zerolog.DebugLevel,
				},
			},
			wantLogger: log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Hook(SeverityHook{}).With().Caller().Timestamp().Logger().Level(zerolog.DebugLevel),
		},
		{
			name: "Scenario 3 - OK - config as nil",
			params: params{
				config: nil,
			},
			wantLogger: zerolog.New(os.Stdout).Hook(SeverityHook{}).With().Caller().Timestamp().Logger().Level(zerolog.InfoLevel),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLogger := New(tt.params.config); !reflect.DeepEqual(gotLogger, tt.wantLogger) {
				t.Errorf("New() = %v, want %v", gotLogger, tt.wantLogger)
			}
		})
	}
}

func TestSeverityHook_Run(t *testing.T) {
	type params struct {
		e     *zerolog.Event
		level zerolog.Level
		msg   string
	}
	tests := []struct {
		name       string
		params     params
		wantResult *zerolog.Event
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				e:     log.Log(),
				level: zerolog.InfoLevel,
			},
			wantResult: (log.Log()).Str("severity", zerolog.InfoLevel.String()),
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				e:     &zerolog.Event{},
				level: zerolog.NoLevel,
			},
			wantResult: &zerolog.Event{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			severityHook := SeverityHook{}

			severityHook.Run(tt.params.e, tt.params.level, tt.params.msg)
			if !reflect.DeepEqual(tt.params.e, tt.wantResult) {
				t.Errorf("Run() = %v, want %v", tt.params.e, tt.wantResult)
			}
		})
	}
}
