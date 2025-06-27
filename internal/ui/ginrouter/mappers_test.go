package ginrouter

import (
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"reflect"
	"testing"
	"time"
)

func Test_fromDomainFnbRequestsToUIFnbRequests(t *testing.T) {
	now := time.Now()
	type params struct {
		from []*models.FnbRequest
	}
	tests := []struct {
		name   string
		params params
		want   []*fnbRequest
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: []*models.FnbRequest{
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
			want: []*fnbRequest{
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
			want: []*fnbRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDomainFnbRequestsToUIFnbRequests(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainFnbRequestsToUIFnbRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDomainFnbRequestToUIFnbRequest(t *testing.T) {

	now := time.Now()
	type params struct {
		from *models.FnbRequest
	}
	tests := []struct {
		name   string
		params params
		want   *fnbRequest
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: &models.FnbRequest{
					ID:          1,
					RequestDate: &now,
					N1:          2,
					S1:          "two",
					N2:          3,
					S2:          "three",
					Limit:       4,
				},
			},
			want: &fnbRequest{
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
				from: &models.FnbRequest{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
			},
			want: &fnbRequest{
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
			if got := fromDomainFnbRequestToUIFnbRequest(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainFnbRequestToUIFnbRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDomainFnbRequestInputStatsToUIFnbRequestInputStats(t *testing.T) {
	type params struct {
		from *models.FnbRequestInputStats
	}
	tests := []struct {
		name   string
		params params
		want   *fnbRequestInputStats
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: &models.FnbRequestInputStats{
					N1:    2,
					S1:    "two",
					N2:    3,
					S2:    "three",
					Limit: 4,
					Count: 5,
				},
			},
			want: &fnbRequestInputStats{
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
				from: &models.FnbRequestInputStats{
					N1:    6,
					S1:    "six",
					N2:    7,
					S2:    "seven",
					Limit: 8,
					Count: 9,
				},
			},
			want: &fnbRequestInputStats{
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
			if got := fromDomainFnbRequestInputStatsToUIFnbRequestInputStats(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainFnbRequestInputStatsToUIFnbRequestInputStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDomainFnbRequestsInputsStatsToUIFnbRequestsInputsStats(t *testing.T) {
	type params struct {
		from []*models.FnbRequestInputStats
	}
	tests := []struct {
		name   string
		params params
		want   []*fnbRequestInputStats
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				from: []*models.FnbRequestInputStats{
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
			want: []*fnbRequestInputStats{
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
			want: []*fnbRequestInputStats{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDomainFnbRequestsInputsStatsToUIFnbRequestsInputsStats(tt.params.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDomainFnbRequestsInputsStatsToUIFnbRequestsInputsStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
