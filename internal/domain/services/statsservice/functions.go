package statsservice

import (
	"fmt"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
)

// GetFnbRequestsInputsStats returns data related to provided arguments to GetFizzAndBuzz function
func (service *statsService) GetFnbRequestsInputsStats(sorted bool) ([]*models.FnbRequestInputStats, error) {
	result, err := service.dal.GetFnbRequestsInputsStats(sorted)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve data to compute for statistics: %w", err)
	}

	return result, nil
}

// GetFnbRequestsMostUsedCombination returns the most provided arguments to GetFizzAndBuzz function
func (service *statsService) GetFnbRequestsMostUsedCombination() (*models.FnbRequestInputStats, error) {
	result, err := service.dal.GetFnbRequestsMostUsedCombination()
	if err != nil {
		return nil, fmt.Errorf("couldn't get most used combination: %w", err)
	}

	return result, nil
}
