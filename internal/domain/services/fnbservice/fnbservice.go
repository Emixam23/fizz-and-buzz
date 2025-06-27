package fnbservice

import (
	"errors"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
)

// Config contains the configuration for the FNB Service
type Config struct {
	Zero uint64
}

// IFnbService is a service that handles input fizz-and-buzz commands from users. It can also return the history of the executed command.
type IFnbService interface {
	GetFizzAndBuzz(m1 uint32, s1 string, m2 uint32, s2 string, limit uint64) ([]string, error)
	GetFnbRequestsHistory(limit *uint64) ([]*models.FnbRequest, error)
}

type fnbService struct {
	dal infra.IDAL

	zero uint64
}

// New creates and initializes a new Fizz And Buzz service
func New(config *Config, dal infra.IDAL) (IFnbService, error) {

	if config == nil {
		return nil, errors.New("no config provided")
	} else if dal == nil {
		return nil, errors.New("provided database access layer does not seem initialized (nil)")
	}

	return &fnbService{
		dal:  dal,
		zero: config.Zero,
	}, nil
}
