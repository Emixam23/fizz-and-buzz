package integration

import (
	"bou.ke/monkey"
	"context"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/emixam23/fizz-and-buzz/internal"
	"gitlab.com/emixam23/fizz-and-buzz/internal/config"
	"gitlab.com/emixam23/fizz-and-buzz/internal/infra/dal"
	"net/http/httptest"
	"os"
)

type TestContext struct {
	app internal.IApp

	conf              *config.Config
	dbClient          *sql.DB
	mockedTime        *monkey.PatchGuard
	responsesRecorder []*httptest.ResponseRecorder
}

func (testContext *TestContext) Before(ctx context.Context, s *godog.Scenario) (context.Context, error) {

	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("an error happened while loading the config: %w", err)
	}

	dbClient, err := dal.NewPostgreSQLDbClient(conf.DatabaseConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't init db client: %w", err)
	}

	app, err := internal.NewWithArgs(conf, dbClient)
	if err != nil {
		return ctx, fmt.Errorf("an error happened while creating/initializing new app with args: %w", err)
	}

	testContext.conf = conf
	testContext.app = app
	testContext.dbClient = dbClient
	testContext.responsesRecorder = make([]*httptest.ResponseRecorder, 0)

	// cleanup/recreate table for tests
	_ = execQueryAgainstDatabase("DROP TABLE IF EXISTS fnb_requests")
	createTableQuery := "CREATE TABLE IF NOT EXISTS fnb_requests (" +
		"id SERIAL PRIMARY KEY, " +
		"request_date TIMESTAMP NOT NULL, " +
		"n1 BIGSERIAL NOT NULL, " +
		"s1 TEXT NOT NULL, " +
		"n2 BIGSERIAL NOT NULL, " +
		"s2 TEXT NOT NULL, " +
		"rlimit BIGSERIAL NOT NULL" +
		");"
	if err := execQueryAgainstDatabase(createTableQuery); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (testContext *TestContext) After(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {

	// mocked time
	if testContext.mockedTime != nil {
		testContext.mockedTime.Unpatch()
	}

	// cleanup database after test
	_ = execQueryAgainstDatabase("DROP TABLE IF EXISTS fnb_requests")

	return ctx, nil
}

func execQueryAgainstDatabase(query string) error {
	tx, err := testContext.dbClient.Begin()
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("couldn't prepare statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("couldn't execute statement: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}

func loadConfig() (*config.Config, error) {

	givenEnv := map[string]string{
		"LOGGER_LEVEL":                  "warn",
		"LOGGER_AS_JSON_FORMAT":         "false",
		"DATABASE_NAME":                 "testdb_integration_tests",
		"DATABASE_USER":                 "testapi_integration_tests",
		"DATABASE_PASSWORD":             "fizznbuzz_integration_tests",
		"DATABASE_TIMEOUT_SECONDS":      "5",
		"DATABASE_RETRY_AMOUNT_ON_FAIL": "7",
		"FNB_SERVICE_ZERO":              "1",
		"ROUTER_HOST":                   "localhost",
		"ROUTER_PORT":                   "9000",
		"ROUTER_MODE":                   "test",
	}

	for key, value := range givenEnv {
		if err := os.Setenv(key, value); err != nil {
			return nil, fmt.Errorf("couldn't set env key/value: %w", err)
		}
	}

	c := config.LoadConfig()
	return c, nil
}
