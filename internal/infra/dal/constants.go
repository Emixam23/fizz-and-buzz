package dal

const (
	selectFnbRequestsSortedDescByIDQuery             = "SELECT id, request_date, n1, s1, n2, s2, rlimit FROM fnb_requests ORDER BY id DESC"
	selectFnbRequestsWithLimitQuery                  = selectFnbRequestsSortedDescByIDQuery + " LIMIT $1"
	selectFnbRequestsInputsStatsQuery                = "SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count FROM fnb_requests GROUP BY n1, s1, n2, s2, rlimit"
	selectFnbRequestsInputsStatsSortedQuery          = selectFnbRequestsInputsStatsQuery + " ORDER BY fnb_count DESC"
	selectFnbRequestsInputsStatsSortedWithLimitQuery = selectFnbRequestsInputsStatsSortedQuery + " LIMIT $1"
	registerFnbRequestInputQuery                     = "INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)"
)
