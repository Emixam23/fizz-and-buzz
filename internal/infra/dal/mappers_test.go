package dal

import (
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"reflect"
	"testing"
	"time"
)

func Test_fromDataModelFnbRequestToDomainFnbRequest(t *testing.T) {
	now := time.Now()
	type params struct {
		from *fnbRequest
	}
	tests := []struct {
		name   string
		params params
		want   *models.FnbRequest
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: &fnbRequest{
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
			want: &models.FnbRequest{
				ID:          1,
				RequestDate: &now,
				N1:          2,
				S1:          "two",
				N2:          3,
				S2:          "three",
				Limit:       4,
			},
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				from: &fnbRequest{
					tableName:   struct{}{},
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
			},
			want: &models.FnbRequest{
				ID:          5,
				RequestDate: &now,
				N1:          6,
				S1:          "six",
				N2:          7,
				S2:          "seven",
				Limit:       8,
			},
		},
		{
			name: "Scenario 3 - OK",
			params: params{
				from: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDataModelFnbRequestToDomainFnbRequest(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDataModelFnbRequestToDomainFnbRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDataModelFnbRequestsToDomainFnbRequests(t *testing.T) {
	now := time.Now()
	type params struct {
		from []*fnbRequest
	}
	tests := []struct {
		name   string
		params params
		want   []*models.FnbRequest
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: []*fnbRequest{
					{
						ID:          1,
						RequestDate: &now,
						N1:          2,
						S1:          "two",
						N2:          3,
						S2:          "three",
						Limit:       4,
					},
					{
						ID:          5,
						RequestDate: &now,
						N1:          6,
						S1:          "six",
						N2:          7,
						S2:          "seven",
						Limit:       8,
					},
				},
			},
			want: []*models.FnbRequest{
				{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
			},
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				from: nil,
			},
			want: []*models.FnbRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDataModelFnbRequestsToDomainFnbRequests(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDataModelFnbRequestsToDomainFnbRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats(t *testing.T) {
	type params struct {
		from *fnbRequestInputStats
	}
	tests := []struct {
		name   string
		params params
		want   *models.FnbRequestInputStats
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: &fnbRequestInputStats{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 5,
				},
			},
			want: &models.FnbRequestInputStats{
				N1:    2,
				S1:    "two",
				N2:    3,
				S2:    "three",
				Limit: 4,
				Count: 5,
			},
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				from: &fnbRequestInputStats{
					N1:    6,
					S1:    "six",
					N2:    7,
					S2:    "seven",
					Limit: 8,
					Count: 9,
				},
			},
			want: &models.FnbRequestInputStats{
				N1:    6,
				S1:    "six",
				N2:    7,
				S2:    "seven",
				Limit: 8,
				Count: 9,
			},
		},
		{
			name: "Scenario 3 - OK",
			params: params{
				from: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDataModelFnbRequestInputStatsToDomainFnbRequestInputStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDataModelFnbRequestsInputsStatsToDomainFnbRequestsInputsStats(t *testing.T) {
	type params struct {
		from []*fnbRequestInputStats
	}
	tests := []struct {
		name   string
		params params
		want   []*models.FnbRequestInputStats
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: []*fnbRequestInputStats{
					{
						N1:    2,
						S1:    "two",
						N2:    3,
						S2:    "three",
						Limit: 4,
						Count: 5,
					},
					{
						N1:    6,
						S1:    "six",
						N2:    7,
						S2:    "seven",
						Limit: 8,
						Count: 9,
					},
				},
			},
			want: []*models.FnbRequestInputStats{
				{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 5,
				},
				{
					N1:    6,
					S1:    "six",
					N2:    7,
					S2:    "seven",
					Limit: 8,
					Count: 9,
				},
			},
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				from: nil,
			},
			want: []*models.FnbRequestInputStats{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDataModelFnbRequestsInputsStatsToDomainFnbRequestsInputsStats(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDataModelFnbRequestsInputsStatsToDomainFnbRequestsInputsStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
