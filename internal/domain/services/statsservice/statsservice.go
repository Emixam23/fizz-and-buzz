package statsservice

import (
	"errors"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
)

// IStatsService returns computed data regarding statistics using the DAL
type IStatsService interface {
	GetFnbRequestsInputsStats(sorted bool) ([]*models.FnbRequestInputStats, error)
	GetFnbRequestsMostUsedCombination() (*models.FnbRequestInputStats, error)
}

type statsService struct {
	dal infra.IDAL
}

// New creates and initializes a new Statistics service
func New(dal infra.IDAL) (IStatsService, error) {

	if dal == nil {
		return nil, errors.New("provided database access layer does not seem initialized (nil)")
	}

	return &statsService{dal: dal}, nil
}
