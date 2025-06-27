package internal

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gitlab.com/emixam23/fizz-and-buzz/internal/config"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/utils/logger"
	"reflect"
	"testing"
)

func TestNewWithArgs(t *testing.T) {

	client, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type args struct {
		conf     *config.Config
		dbClient *sql.DB
	}
	type wantResult struct {
		appNotNil             bool
		appConfigNotNil       bool
		appDbClientNotNil     bool
		appDalNotNil          bool
		appFnbServiceNotNil   bool
		appStatsServiceNotNil bool
		appRouterNotNil       bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult wantResult
		wantErr    error
	}{
		{
			name: "Scenario 1 - OK",
			args: args{
				conf: &config.Config{
					Logger: &logger.Config{
						JSON:  true,
						Level: zerolog.InfoLevel,
					},
					DatabaseConfig: &infra.Config{
						Host:                  "0.0.0.0",
						Port:                  1234,
						DatabaseName:          "testname",
						User:                  "testuser",
						Password:              "testpassword",
						DefaultTimeoutSeconds: 1,
						RetryAmountOnFail:     1,
					},
					FnbServiceConfig: &fnbservice.Config{Zero: 1},
					RestAPIConfig: &ui.Config{
						Host: "0.0.0.0",
						Port: 9000,
						Mode: "release",
					},
				},
				dbClient: client,
			},
			wantResult: wantResult{
				appNotNil:             true,
				appConfigNotNil:       true,
				appDbClientNotNil:     true,
				appDalNotNil:          true,
				appFnbServiceNotNil:   true,
				appStatsServiceNotNil: true,
				appRouterNotNil:       true,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWithArgs(tt.args.conf, tt.args.dbClient)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewWithArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResult.appNotNil {
				assert.NotNil(t, got, "app should not be nil")
			} else {
				assert.Nil(t, got, "app should be nil")
			}
			if tt.wantResult.appConfigNotNil {
				assert.NotNil(t, got.(*app).config, "app config should not be nil")
			} else {
				assert.Nil(t, got.(*app).config, "app config should be nil")
			}
			if tt.wantResult.appDbClientNotNil {
				assert.NotNil(t, got.(*app).dbClient, "app dbClient should not be nil")
			} else {
				assert.Nil(t, got.(*app).dbClient, "app dbClient should be nil")
			}
			if tt.wantResult.appDalNotNil {
				assert.NotNil(t, got.(*app).dal, "app dal should not be nil")
			} else {
				assert.Nil(t, got.(*app).dal, "app dal should be nil")
			}
			if tt.wantResult.appFnbServiceNotNil {
				assert.NotNil(t, got.(*app).fnbService, "app fnb service should not be nil")
			} else {
				assert.Nil(t, got.(*app).fnbService, "app fnb service should be nil")
			}
			if tt.wantResult.appStatsServiceNotNil {
				assert.NotNil(t, got.(*app).statsService, "app stats service should not be nil")
			} else {
				assert.Nil(t, got.(*app).statsService, "app stats service should be nil")
			}
			if tt.wantResult.appRouterNotNil {
				assert.NotNil(t, got.(*app).router, "app router should not be nil")
			} else {
				assert.Nil(t, got.(*app).router, "app router should be nil")
			}
		})
	}
}
