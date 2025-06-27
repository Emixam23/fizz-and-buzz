package dal

import (
	"bou.ke/monkey"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"reflect"
	"testing"
	"time"
)

func Test_dal_GetFnbRequestsHistory(t *testing.T) {

	now := time.Now()
	patchGuard := monkey.Patch(time.Now, func() time.Time { return now })
	defer patchGuard.Unpatch()

	client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	var limit uint64 = 1

	type params struct {
		limit *uint64
	}
	type interactionErrs struct {
		queryErr   error
		rowScanErr bool
	}
	tests := []struct {
		name                 string
		client               *sql.DB
		params               params
		interactionResponses []*fnbRequest
		interactionErrs      interactionErrs
		wantResult           []*models.FnbRequest
		wantErr              error
	}{
		{
			name:   "Scenario 1 - OK",
			client: client,
			params: params{
				limit: &limit,
			},
			interactionResponses: []*fnbRequest{
				{
					tableName:   struct{}{},
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
			},
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: false,
			},
			wantResult: []*models.FnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
			},
			wantErr: nil,
		},
		{
			name:   "Scenario 2 - KO - query error",
			client: client,
			params: params{
				limit: nil,
			},
			interactionResponses: nil,
			interactionErrs: interactionErrs{
				queryErr:   errors.New("any error"),
				rowScanErr: false,
			},
			wantResult: nil,
			wantErr:    fmt.Errorf("couldn't query postgresql database: %w", errors.New("any error")),
		},
		{
			name:   "Scenario 3 - KO - row scan error",
			client: client,
			params: params{
				limit: nil,
			},
			interactionResponses: nil,
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: true,
			},
			wantResult: nil,
			wantErr:    fmt.Errorf("couldn't scan postgresql database result: %w", errors.New("sql: Scan error on column index 0, name \"id\": converting driver.Value type int64 (\"-1\") to a uint64: invalid syntax")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rows := mock.NewRowsWithColumnDefinition(
				mock.NewColumn("id"),
				mock.NewColumn("request_date"),
				mock.NewColumn("n1"),
				mock.NewColumn("s1"),
				mock.NewColumn("n2"),
				mock.NewColumn("s2"),
				mock.NewColumn("rlimit"))

			for _, interactionResponse := range tt.interactionResponses {
				rows.AddRow(interactionResponse.ID, interactionResponse.RequestDate, interactionResponse.N1, interactionResponse.S1, interactionResponse.N2, interactionResponse.S2, interactionResponse.Limit)
			}

			var query *sqlmock.ExpectedQuery
			if tt.params.limit != nil {
				query = mock.ExpectQuery("SELECT id, request_date, n1, s1, n2, s2, rlimit FROM fnb_requests ORDER BY id DESC LIMIT $1").WillReturnRows(rows)
			} else {
				query = mock.ExpectQuery("SELECT id, request_date, n1, s1, n2, s2, rlimit FROM fnb_requests ORDER BY id DESC").WillReturnRows(rows)
			}

			if tt.wantErr != nil {
				if tt.interactionErrs.queryErr != nil {
					query.WillReturnError(tt.interactionErrs.queryErr)
				} else if tt.interactionErrs.rowScanErr {
					rows.AddRow(-1, now, 2, "two", 3, "three", 4) // id as int, not uint
				}
			}

			dal := &dal{client: tt.client}

			got, err := dal.GetFnbRequestsHistory(tt.params.limit)
			if tt.interactionErrs.rowScanErr {
				if err.Error() != tt.wantErr.Error() {
					// to avoid useless complexity on sub error models mocking
					t.Errorf("GetFnbRequestsHistory() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("GetFnbRequestsHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("GetFnbRequestsHistory() got = %v, want %v", got, tt.wantResult)
				return
			}
		})
	}
}

func Test_dal_RegisterFnbRequest(t *testing.T) {

	now := time.Now()
	patchGuard := monkey.Patch(time.Now, func() time.Time { return now })
	defer patchGuard.Unpatch()

	client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type params struct {
		n1    uint32
		s1    string
		n2    uint32
		s2    string
		limit uint64
	}
	tests := []struct {
		name    string
		client  *sql.DB
		params  params
		wantErr error
	}{
		{
			name:   "Scenario 1 - OK",
			client: client,
			params: params{
				n1:    1,
				s1:    "one",
				n2:    2,
				s2:    "two",
				limit: 3,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dal := &dal{client: tt.client}

			mock.ExpectBegin()
			mock.ExpectPrepare("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)")
			mock.ExpectExec("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WithArgs(now, tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			if err := dal.RegisterFnbRequest(tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("RegisterFnbRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dal_RegisterFnbRequest_with_errors_on_interaction(t *testing.T) {

	now := time.Now()
	patchGuard := monkey.Patch(time.Now, func() time.Time { return now })
	defer patchGuard.Unpatch()

	client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type params struct {
		n1    uint32
		s1    string
		n2    uint32
		s2    string
		limit uint64
	}
	type interactionErrs struct {
		beginErr   error
		prepareErr error
		execErr    error
		commitErr  error
	}
	tests := []struct {
		name            string
		client          *sql.DB
		params          params
		interactionErrs interactionErrs
		wantErr         error
	}{
		{
			name:   "Scenario 1 - KO - begin returns an error",
			client: client,
			params: params{
				n1:    1,
				s1:    "one",
				n2:    2,
				s2:    "two",
				limit: 3,
			},
			interactionErrs: interactionErrs{
				beginErr:   errors.New("any error"),
				prepareErr: nil,
				execErr:    nil,
				commitErr:  nil,
			},
			wantErr: fmt.Errorf("couldn't begin transaction: %w", errors.New("any error")),
		},
		{
			name:   "Scenario 1 - KO - prepare returns an error",
			client: client,
			params: params{
				n1:    1,
				s1:    "one",
				n2:    2,
				s2:    "two",
				limit: 3,
			},
			interactionErrs: interactionErrs{
				beginErr:   nil,
				prepareErr: errors.New("any error"),
				execErr:    nil,
				commitErr:  nil,
			},
			wantErr: fmt.Errorf("couldn't prepare statement: %w", errors.New("any error")),
		},
		{
			name:   "Scenario 1 - KO - exec returns an error",
			client: client,
			params: params{
				n1:    1,
				s1:    "one",
				n2:    2,
				s2:    "two",
				limit: 3,
			},
			interactionErrs: interactionErrs{
				beginErr:   nil,
				prepareErr: nil,
				execErr:    errors.New("any error"),
				commitErr:  nil,
			},
			wantErr: fmt.Errorf("couldn't execute statement: %w", errors.New("any error")),
		},
		{
			name:   "Scenario 1 - KO - commit returns an error",
			client: client,
			params: params{
				n1:    1,
				s1:    "one",
				n2:    2,
				s2:    "two",
				limit: 3,
			},
			interactionErrs: interactionErrs{
				beginErr:   nil,
				prepareErr: nil,
				execErr:    nil,
				commitErr:  errors.New("any error"),
			},
			wantErr: fmt.Errorf("couldn't commit transaction: %w", errors.New("any error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dal := &dal{client: tt.client}

			if tt.interactionErrs.beginErr != nil {
				mock.ExpectBegin().WillReturnError(tt.interactionErrs.beginErr)
			} else if tt.interactionErrs.prepareErr != nil {
				mock.ExpectBegin().WillReturnError(tt.interactionErrs.beginErr)
				mock.ExpectPrepare("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WillReturnError(tt.interactionErrs.prepareErr)
			} else if tt.interactionErrs.execErr != nil {
				mock.ExpectBegin().WillReturnError(tt.interactionErrs.beginErr)
				mock.ExpectPrepare("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WillReturnError(tt.interactionErrs.prepareErr)
				mock.ExpectExec("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WithArgs(now, tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.interactionErrs.execErr)
			} else if tt.interactionErrs.commitErr != nil {
				mock.ExpectBegin().WillReturnError(tt.interactionErrs.beginErr)
				mock.ExpectPrepare("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WillReturnError(tt.interactionErrs.prepareErr)
				mock.ExpectExec("INSERT INTO fnb_requests (request_date, n1, s1, n2, s2, rlimit) VALUES ($1, $2, $3, $4, $5, $6)").WithArgs(now, tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.interactionErrs.execErr)
				mock.ExpectCommit().WillReturnError(tt.interactionErrs.commitErr)
			}

			if err := dal.RegisterFnbRequest(tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("RegisterFnbRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
