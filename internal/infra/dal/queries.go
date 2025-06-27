package dal

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"time"
)

func (dal *dal) GetFnbRequestsHistory(limit *uint64) ([]*models.FnbRequest, error) {

	var rows *sql.Rows
	var err error

	if limit != nil {
		rows, err = dal.client.Query(selectFnbRequestsWithLimitQuery, limit)
	} else {
		rows, err = dal.client.Query(selectFnbRequestsSortedDescByIDQuery)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't query postgresql database: %w", err)
	}
	defer rows.Close()

	fnbRequests := make([]*fnbRequest, 0)
	for rows.Next() {
		var req fnbRequest
		if err := rows.Scan(&req.ID, &req.RequestDate, &req.N1, &req.S1, &req.N2, &req.S2, &req.Limit); err != nil {
			return nil, fmt.Errorf("couldn't scan postgresql database result: %w", err)
		}
		fnbRequests = append(fnbRequests, &req)
	}
	return fromDataModelFnbRequestsToDomainFnbRequests(fnbRequests), nil
}

func (dal *dal) RegisterFnbRequest(n1 uint32, s1 string, n2 uint32, s2 string, limit uint64) error {

	now := time.Now()

	tx, err := dal.client.Begin()
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(registerFnbRequestInputQuery)
	if err != nil {
		return fmt.Errorf("couldn't prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(now, n1, s1, n2, s2, limit)
	if err != nil {
		return fmt.Errorf("couldn't execute statement: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	log.Trace().
		Time("request_date", now).
		Uint32("n1", n1).
		Str("s1", s1).
		Uint32("n2", n2).
		Str("s2", s2).
		Uint64("limit", limit).
		Msg("Inserted fizz and buzz request into postgresql")

	return nil
}
