package infra

import (
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
)

// IDAL is the Data Access Layer of our app
type IDAL interface {
	RegisterFnbRequest(n1 uint32, s1 string, n2 uint32, s2 string, limit uint64) error
	GetFnbRequestsHistory(limit *uint64) ([]*models.FnbRequest, error)

	GetFnbRequestsInputsStats(sorted bool) ([]*models.FnbRequestInputStats, error)
	GetFnbRequestsMostUsedCombination() (*models.FnbRequestInputStats, error)
}

// Config contains the configuration for the DAL
type Config struct {
	Host                         string
	Port                         uint64
	DatabaseName, User, Password string
	DefaultTimeoutSeconds        uint
	RetryAmountOnFail            uint32
}
