package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"reflect"
	"testing"
)

func Test_dal_GetFnbRequestsInputsStats(t *testing.T) {

	client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type params struct {
		sorted bool
	}
	type interactionErrs struct {
		queryErr   error
		rowScanErr bool
	}
	tests := []struct {
		name                 string
		client               *sql.DB
		params               params
		interactionResponses []*fnbRequestInputStats
		interactionErrs      interactionErrs
		wantResult           []*models.FnbRequestInputStats
		wantErr              error
	}{
		{
			name:   "Scenario 1 - OK",
			client: client,
			params: params{
				sorted: false,
			},
			interactionResponses: []*fnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 1,
				},
			},
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: false,
			},
			wantResult: []*models.FnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 1,
				},
			},
			wantErr: nil,
		},
		{
			name:   "Scenario 2 - OK - sorted",
			client: client,
			params: params{
				sorted: true,
			},
			interactionResponses: []*fnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 1,
				},
			},
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: false,
			},
			wantResult: []*models.FnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 1,
				},
			},
			wantErr: nil,
		},
		{
			name:   "Scenario 3 - KO - query error",
			client: client,
			params: params{
				sorted: false,
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
			name:   "Scenario 4 - KO - row scan error",
			client: client,
			params: params{
				sorted: false,
			},
			interactionResponses: nil,
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: true,
			},
			wantResult: nil,
			wantErr:    fmt.Errorf("couldn't scan postgresql database result: %w", errors.New("sql: Scan error on column index 5, name \"fnb_count\": converting driver.Value type string (\"missmatch type\") to a uint64: invalid syntax")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rows := mock.NewRowsWithColumnDefinition(
				mock.NewColumn("n1"),
				mock.NewColumn("s1"),
				mock.NewColumn("n2"),
				mock.NewColumn("s2"),
				mock.NewColumn("rlimit"),
				mock.NewColumn("fnb_count"))

			for _, interactionResponse := range tt.interactionResponses {
				rows.AddRow(interactionResponse.N1, interactionResponse.S1, interactionResponse.N2, interactionResponse.S2, interactionResponse.Limit, interactionResponse.Count)
			}

			var query *sqlmock.ExpectedQuery
			if tt.params.sorted {
				query = mock.ExpectQuery("SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count FROM fnb_requests GROUP BY n1, s1, n2, s2, rlimit ORDER BY fnb_count DESC").WillReturnRows(rows)
			} else {
				query = mock.ExpectQuery("SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count FROM fnb_requests GROUP BY n1, s1, n2, s2, rlimit").WillReturnRows(rows)
			}

			if tt.wantErr != nil {
				if tt.interactionErrs.queryErr != nil {
					query.WillReturnError(tt.interactionErrs.queryErr)
				} else if tt.interactionErrs.rowScanErr {
					rows.AddRow(2, "two", 3, "three", 4, "missmatch type") // id as int, not uint
				}
			}

			dal := &dal{client: tt.client}

			got, err := dal.GetFnbRequestsInputsStats(tt.params.sorted)
			if tt.interactionErrs.rowScanErr {
				if err.Error() != tt.wantErr.Error() {
					// to avoid useless complexity on sub error models mocking
					t.Errorf("GetFnbRequestsInputsStats() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("GetFnbRequestsInputsStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("GetFnbRequestsInputsStats() got = %v, want %v", got, tt.wantResult)
				return
			}
		})
	}
}

func Test_dal_GetFnbRequestsMostUsedCombination(t *testing.T) {
	client, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type interactionErrs struct {
		queryErr   error
		rowScanErr bool
	}
	tests := []struct {
		name                 string
		client               *sql.DB
		interactionResponses []*fnbRequestInputStats
		interactionErrs      interactionErrs
		wantResult           *models.FnbRequestInputStats
		wantErr              error
	}{
		{
			name:   "Scenario 1 - OK",
			client: client,
			interactionResponses: []*fnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 1,
				},
			},
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: false,
			},
			wantResult: &models.FnbRequestInputStats{
				N1:    2,
				S1:    "two",
				N2:    3,
				S2:    "three",
				Limit: 4,
				Count: 1,
			},
			wantErr: nil,
		},
		{
			name:                 "Scenario 2 - KO - no data retrieved",
			client:               client,
			interactionResponses: []*fnbRequestInputStats{},
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: false,
			},
			wantResult: nil,
			wantErr:    nil,
		},
		{
			name:                 "Scenario 3 - KO - query error",
			client:               client,
			interactionResponses: nil,
			interactionErrs: interactionErrs{
				queryErr:   errors.New("any error"),
				rowScanErr: false,
			},
			wantResult: nil,
			wantErr:    fmt.Errorf("couldn't query postgresql database: %w", errors.New("any error")),
		},
		{
			name:                 "Scenario 4 - KO - row scan error",
			client:               client,
			interactionResponses: nil,
			interactionErrs: interactionErrs{
				queryErr:   nil,
				rowScanErr: true,
			},
			wantResult: nil,
			wantErr:    fmt.Errorf("couldn't scan postgresql database result: %w", errors.New("sql: Scan error on column index 5, name \"fnb_count\": converting driver.Value type string (\"missmatch type\") to a uint64: invalid syntax")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rows := mock.NewRowsWithColumnDefinition(
				mock.NewColumn("n1"),
				mock.NewColumn("s1"),
				mock.NewColumn("n2"),
				mock.NewColumn("s2"),
				mock.NewColumn("rlimit"),
				mock.NewColumn("fnb_count"))

			for _, interactionResponse := range tt.interactionResponses {
				rows.AddRow(interactionResponse.N1, interactionResponse.S1, interactionResponse.N2, interactionResponse.S2, interactionResponse.Limit, interactionResponse.Count)
			}

			query := mock.ExpectQuery("SELECT n1, s1, n2, s2, rlimit, COUNT(*) AS fnb_count FROM fnb_requests GROUP BY n1, s1, n2, s2, rlimit ORDER BY fnb_count DESC LIMIT $1").WithArgs(1).WillReturnRows(rows)

			if tt.wantErr != nil {
				if tt.interactionErrs.queryErr != nil {
					query.WillReturnError(tt.interactionErrs.queryErr)
				} else if tt.interactionErrs.rowScanErr {
					rows.AddRow(2, "two", 3, "three", 4, "missmatch type") // id as int, not uint
				}
			}

			dal := &dal{client: tt.client}

			got, err := dal.GetFnbRequestsMostUsedCombination()
			if tt.interactionErrs.rowScanErr {
				if err.Error() != tt.wantErr.Error() {
					// to avoid useless complexity on sub error models mocking
					t.Errorf("GetFnbRequestsMostUsedCombination() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("GetFnbRequestsMostUsedCombination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("GetFnbRequestsMostUsedCombination() got = %v, want %v", got, tt.wantResult)
				return
			}
		})
	}
}
