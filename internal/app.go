package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"gitlab.com/emixam23/fizz-and-buzz/internal/config"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/fnbservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/services/statsservice"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/ui"
	"gitlab.com/emixam23/fizz-and-buzz/internal/infra/dal"
	"gitlab.com/emixam23/fizz-and-buzz/internal/ui/ginrouter"
)

// IApp contains the whole app and is handling its running status
type IApp interface {
	Start() error

	// GetRouter is exposed for integration testing requirement
	GetRouter() ui.IRestAPI
}

type app struct {
	config *config.Config

	dbClient *sql.DB
	dal      infra.IDAL

	fnbService   fnbservice.IFnbService
	statsService statsservice.IStatsService

	router ui.IRestAPI
}

// New creates a new IApp fully initialized and ready to run
func New() (IApp, error) {

	conf := config.LoadConfig()
	dbClient, err := dal.NewPostgreSQLDbClient(conf.DatabaseConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't init db client: %w", err)
	}

	return NewWithArgs(conf, dbClient)
}

// NewWithArgs creates a new IApp fully initialized and ready to run
// It (unlike New) requires a preloaded configuration and a db client
func NewWithArgs(conf *config.Config, dbClient *sql.DB) (IApp, error) {

	if conf == nil {
		return nil, errors.New("provided config is nil")
	} else if dbClient == nil {
		return nil, errors.New("provided db client is nil")
	}

	app := &app{
		config:   conf,
		dbClient: dbClient,
	}
	var err error

	// infra
	if app.dal, err = dal.New(app.dbClient); err != nil {
		return nil, fmt.Errorf("couldn't init dal: %w", err)
	}

	// domain
	if app.fnbService, err = fnbservice.New(conf.FnbServiceConfig, app.dal); err != nil {
		return nil, fmt.Errorf("couldn't init fnb service: %w", err)
	}
	if app.statsService, err = statsservice.New(app.dal); err != nil {
		return nil, fmt.Errorf("couldn't init stats service: %w", err)
	}

	// ui
	if app.router, err = ginrouter.New(conf.RestAPIConfig, app.fnbService, app.statsService); err != nil {
		return nil, fmt.Errorf("couldn't init router/rest api: %w", err)
	}

	return app, nil
}

func (app *app) Start() error {
	return app.router.ListenAndServe()
}

func (app app) GetRouter() ui.IRestAPI {
	return app.router
}
