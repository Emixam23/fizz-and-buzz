package dal

import (
	"database/sql"
	"fmt"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
)

/*
SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count
FROM fnb_requests
GROUP BY n1, s1, n2, s2, rlimit
ORDER BY fnb_count DESC;
*/

func (dal *dal) GetFnbRequestsInputsStats(sorted bool) ([]*models.FnbRequestInputStats, error) {

	var rows *sql.Rows
	var err error

	if sorted {
		rows, err = dal.client.Query(selectFnbRequestsInputsStatsSortedQuery)
	} else {
		rows, err = dal.client.Query(selectFnbRequestsInputsStatsQuery)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't query postgresql database: %w", err)
	}
	defer rows.Close()

	fnbRequestsInputsStats := make([]*fnbRequestInputStats, 0)
	for rows.Next() {
		var stats fnbRequestInputStats
		if err := rows.Scan(&stats.N1, &stats.S1, &stats.N2, &stats.S2, &stats.Limit, &stats.Count); err != nil {
			return nil, fmt.Errorf("couldn't scan postgresql database result: %w", err)
		}
		fnbRequestsInputsStats = append(fnbRequestsInputsStats, &stats)
	}
	return fromDataModelFnbRequestsInputsStatsToDomainFnbRequestsInputsStats(fnbRequestsInputsStats), nil
}

/*
SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count
FROM fnb_requests
GROUP BY n1, s1, n2, s2, rlimit
ORDER BY fnb_count DESC
LIMIT 1;
*/

func (dal *dal) GetFnbRequestsMostUsedCombination() (*models.FnbRequestInputStats, error) {

	rows, err := dal.client.Query(selectFnbRequestsInputsStatsSortedWithLimitQuery, 1)

	if err != nil {
		return nil, fmt.Errorf("couldn't query postgresql database: %w", err)
	}
	defer rows.Close()

	fnbRequestsInputsStats := make([]*fnbRequestInputStats, 0)
	for rows.Next() {
		var stats fnbRequestInputStats
		if err := rows.Scan(&stats.N1, &stats.S1, &stats.N2, &stats.S2, &stats.Limit, &stats.Count); err != nil {
			return nil, fmt.Errorf("couldn't scan postgresql database result: %w", err)
		}
		fnbRequestsInputsStats = append(fnbRequestsInputsStats, &stats)
	}

	if len(fnbRequestsInputsStats) == 0 {
		return nil, nil
	}

	return fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats(fnbRequestsInputsStats[0]), nil
}
