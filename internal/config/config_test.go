package config

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/utils/logger"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		givenEnv map[string]string
		want     *Config
	}{
		{
			name: "Scenario 1 - OK - default",
			want: &Config{
				Logger: &logger.Config{
					JSON:  true,
					Level: zerolog.InfoLevel,
				},
				DatabaseConfig: &infra.Config{
					DefaultTimeoutSeconds: 2,
				},
				FnbServiceConfig: &fnbservice.Config{Zero: 1},
				RestAPIConfig: &ui.Config{
					Host: "0.0.0.0",
					Port: 8080,
					Mode: "release",
				},
			},
		},
		{
			name: "Scenario 2 - OK - with full values provided",
			givenEnv: map[string]string{
				"LOGGER_LEVEL":                  "warn",
				"LOGGER_AS_JSON_FORMAT":         "false",
				"DATABASE_HOST":                 "azerty_host",
				"DATABASE_NAME":                 "azerty_name",
				"DATABASE_USER":                 "azerty_user",
				"DATABASE_PASSWORD":             "azerty_pass",
				"DATABASE_PORT":                 "1234",
				"DATABASE_TIMEOUT_SECONDS":      "5",
				"DATABASE_RETRY_AMOUNT_ON_FAIL": "7",
				"FNB_SERVICE_ZERO":              "10",
				"ROUTER_HOST":                   "localhost",
				"ROUTER_PORT":                   "9000",
				"ROUTER_MODE":                   "test",
			},
			want: &Config{
				Logger: &logger.Config{
					JSON:  false,
					Level: zerolog.WarnLevel,
				},
				DatabaseConfig: &infra.Config{
					Host:                  "azerty_host",
					Port:                  1234,
					DatabaseName:          "azerty_name",
					User:                  "azerty_user",
					Password:              "azerty_pass",
					DefaultTimeoutSeconds: 5,
					RetryAmountOnFail:     7,
				},
				FnbServiceConfig: &fnbservice.Config{
					Zero: 10,
				},
				RestAPIConfig: &ui.Config{
					Host: "localhost",
					Port: 9000,
					Mode: "test",
				},
			},
		},
		{
			name: "Scenario 3 - OK - zerolog level is invalid, defaulting to info",
			givenEnv: map[string]string{
				"LOGGER_LEVEL": "azerty",
			},
			want: &Config{
				Logger: &logger.Config{
					JSON:  true,
					Level: zerolog.InfoLevel,
				},
				DatabaseConfig: &infra.Config{
					DefaultTimeoutSeconds: 2,
				},
				FnbServiceConfig: &fnbservice.Config{Zero: 1},
				RestAPIConfig: &ui.Config{
					Host: "0.0.0.0",
					Port: 8080,
					Mode: "release",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.givenEnv != nil {
				for key, value := range tt.givenEnv {
					_ = os.Setenv(key, value)
				}
			}

			got := LoadConfig()
			assert.Equal(t, tt.want, got)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadConfig() got = %v, want %v", got, tt.want)
			//}

			os.Clearenv()
		})
	}
}
