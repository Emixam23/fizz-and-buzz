package ginrouter

import "gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"

/* domain -> ui */

func fromDomainFnbRequestToUIFnbRequest(from *models.FnbRequest) *fnbRequest {
	if from == nil {
		return nil
	}

	to := &fnbRequest{}

	to.ID = from.ID
	to.RequestDate = from.RequestDate
	to.N1 = from.N1
	to.S1 = from.S1
	to.N2 = from.N2
	to.S2 = from.S2
	to.Limit = from.Limit

	return to
}

func fromDomainFnbRequestsToUIFnbRequests(from []*models.FnbRequest) []*fnbRequest {
	var to = make([]*fnbRequest, 0)

	if from == nil {
		return to
	}

	for _, item := range from {
		to = append(to, fromDomainFnbRequestToUIFnbRequest(item))
	}

	return to
}

func fromDomainFnbRequestInputStatsToUIFnbRequestInputStats(from *models.FnbRequestInputStats) *fnbRequestInputStats {
	if from == nil {
		return nil
	}

	to := &fnbRequestInputStats{}

	to.N1 = from.N1
	to.S1 = from.S1
	to.N2 = from.N2
	to.S2 = from.S2
	to.Limit = from.Limit
	to.Count = from.Count

	return to
}

func fromDomainFnbRequestsInputsStatsToUIFnbRequestsInputsStats(from []*models.FnbRequestInputStats) []*fnbRequestInputStats {
	var to = make([]*fnbRequestInputStats, 0)

	if from == nil {
		return to
	}

	for _, item := range from {
		to = append(to, fromDomainFnbRequestInputStatsToUIFnbRequestInputStats(item))
	}

	return to
}
