package dal

import (
	"database/sql"
	"errors"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
)

type dal struct {
	client *sql.DB
}

// New takes a db client as argument and returns a freshly created DAL initialized
func New(client *sql.DB) (infra.IDAL, error) {

	if client == nil {
		return nil, errors.New("provided postgres client is nil")
	}

	return &dal{
		client: client,
	}, nil
}
