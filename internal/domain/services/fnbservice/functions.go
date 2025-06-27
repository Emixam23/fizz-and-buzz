package fnbservice

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models/errormodels"
	"strconv"
)

// GetFizzAndBuzz returns the results for the provided arguments and persist these arguments into a given storage throughout the DAL
func (service *fnbService) GetFizzAndBuzz(n1 uint32, s1 string, n2 uint32, s2 string, limit uint64) ([]string, error) {

	if limit <= service.zero {
		return nil, errormodels.NewUnprocessableError(fmt.Errorf("provided limit \"%d\" must be greater than \"%d\" (provided in configuration)", limit, service.zero))
	}

	// preparing value to avoid duplicate computation/cast
	uint64N1 := uint64(n1)
	uint64N2 := uint64(n2)
	combinedMultiple := uint64N1 * uint64N2

	var results []string
	for i := service.zero; i <= limit; i++ {
		if i%combinedMultiple == 0 {
			results = append(results, s1+s2)
		} else if i%uint64N1 == 0 {
			results = append(results, s1)
		} else if i%uint64N2 == 0 {
			results = append(results, s2)
		} else {
			results = append(results, strconv.FormatUint(i, 10))
		}
	}

	if err := service.dal.RegisterFnbRequest(n1, s1, n2, s2, limit); err != nil {
		log.Warn().
			Err(err).
			Uint32("n1", n1).
			Str("s1", s1).
			Uint32("n2", n2).
			Str("s2", s2).
			Uint64("limit", limit).
			Msg("Couldn't register the processed fizz and buzz request to the storage")
	}

	return results, nil

}

// GetFnbRequestsHistory returns the history of provided arguments (for GetFizzAndBuzz function)
func (service *fnbService) GetFnbRequestsHistory(limit *uint64) ([]*models.FnbRequest, error) {

	fnbRequestsHistory, err := service.dal.GetFnbRequestsHistory(limit)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve fizz and buzz requests history: %w", err)
	}

	return fnbRequestsHistory, nil
}
