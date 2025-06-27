package dal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/utils/retry"
	"time"

	// required to register postgres as driver -> init()
	_ "github.com/lib/pq"
)

// NewPostgreSQLDbClient takes a configuration as argument and create/initializes a new sql client for postgres.
// The function also "ping" the database before returning the client (to ensure the connection is working)
func NewPostgreSQLDbClient(config *infra.Config) (*sql.DB, error) {

	if config == nil {
		return nil, errors.New("provided config is nil")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName)

	client, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("an error happened while opening sql connection: %w", err)
	}

	if client == nil {
		return nil, errors.New("created postgresql instance client is nil")
	}

	if err := retry.Retry(int(config.RetryAmountOnFail), time.Duration(config.DefaultTimeoutSeconds)*time.Second, func(attempts int, sleep time.Duration) error {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.DefaultTimeoutSeconds)*time.Second)
		err := client.PingContext(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Database returned an error on ping")
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("postgresql seems unreachable: %w", err)
	}

	return client, nil
}
