package dal

import "gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"

/* infra -> domain */

func fromDataModelFnbRequestToDomainFnbRequest(from *fnbRequest) *models.FnbRequest {
	if from == nil {
		return nil
	}

	to := &models.FnbRequest{}

	to.ID = from.ID
	to.RequestDate = from.RequestDate
	to.N1 = from.N1
	to.S1 = from.S1
	to.N2 = from.N2
	to.S2 = from.S2
	to.Limit = from.Limit

	return to
}

func fromDataModelFnbRequestsToDomainFnbRequests(from []*fnbRequest) []*models.FnbRequest {
	var to = make([]*models.FnbRequest, 0)

	if from == nil {
		return to
	}

	for _, item := range from {
		to = append(to, fromDataModelFnbRequestToDomainFnbRequest(item))
	}

	return to
}

func fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats(from *fnbRequestInputStats) *models.FnbRequestInputStats {
	if from == nil {
		return nil
	}

	to := &models.FnbRequestInputStats{}

	to.N1 = from.N1
	to.S1 = from.S1
	to.N2 = from.N2
	to.S2 = from.S2
	to.Limit = from.Limit
	to.Count = from.Count

	return to
}

func fromDataModelFnbRequestsInputsStatsToDomainFnbRequestsInputsStats(from []*fnbRequestInputStats) []*models.FnbRequestInputStats {
	var to = make([]*models.FnbRequestInputStats, 0)

	if from == nil {
		return to
	}

	for _, item := range from {
		to = append(to, fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats(item))
	}

	return to
}
