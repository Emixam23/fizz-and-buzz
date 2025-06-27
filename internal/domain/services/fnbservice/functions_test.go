package fnbservice

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models/errormodels"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"reflect"
	"testing"
	"time"
)

func Test_fnbService_GetFizzAndBuzz(t *testing.T) {

	type params struct {
		n1    uint32
		s1    string
		n2    uint32
		s2    string
		limit uint64

		dalInteractionsAmount int
		dalErrorOnInteraction error
	}
	tests := []struct {
		name       string
		config     *Config
		params     params
		wantResult []string
		wantErr    error
	}{
		{
			name:   "Scenario 1 - OK",
			config: &Config{Zero: 1},
			params: params{
				n1:                    3,
				s1:                    "fizz",
				n2:                    5,
				s2:                    "buzz",
				limit:                 19,
				dalInteractionsAmount: 1,
				dalErrorOnInteraction: nil,
			},
			wantResult: []string{
				"1",
				"2",
				"fizz",
				"4",
				"buzz",
				"fizz",
				"7",
				"8",
				"fizz",
				"buzz",
				"11",
				"fizz",
				"13",
				"14",
				"fizzbuzz",
				"16",
				"17",
				"fizz",
				"19",
			},
			wantErr: nil,
		},
		{
			name:   "Scenario 2 - KO  - limit value below service zero value",
			config: &Config{Zero: 25},
			params: params{
				n1:                    3,
				s1:                    "fizz",
				n2:                    5,
				s2:                    "buzz",
				limit:                 19,
				dalInteractionsAmount: 0,
				dalErrorOnInteraction: nil,
			},
			wantResult: nil,
			wantErr:    errormodels.NewUnprocessableError(errors.New("provided limit \"19\" must be greater than \"25\" (provided in configuration)")),
		},
		{
			name:   "Scenario 3 - KO  - database returns an error so a warn log should be displayed in the console",
			config: &Config{Zero: 1},
			params: params{
				n1:                    3,
				s1:                    "fizz",
				n2:                    5,
				s2:                    "buzz",
				limit:                 19,
				dalInteractionsAmount: 1,
				dalErrorOnInteraction: errors.New("any error"),
			},
			wantResult: []string{
				"1",
				"2",
				"fizz",
				"4",
				"buzz",
				"fizz",
				"7",
				"8",
				"fizz",
				"buzz",
				"11",
				"fizz",
				"13",
				"14",
				"fizzbuzz",
				"16",
				"17",
				"fizz",
				"19",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			dalMock := mocks.NewMockIDAL(ctrl)

			if tt.params.dalErrorOnInteraction == nil {
				dalMock.EXPECT().RegisterFnbRequest(tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit).Times(tt.params.dalInteractionsAmount)
			} else {
				dalMock.EXPECT().RegisterFnbRequest(tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit).Times(tt.params.dalInteractionsAmount).Return(errors.New("any error"))
			}

			service := &fnbService{dal: dalMock, zero: tt.config.Zero}

			result, err := service.GetFizzAndBuzz(tt.params.n1, tt.params.s1, tt.params.n2, tt.params.s2, tt.params.limit)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("GetFizzAndBuzz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("GetFizzAndBuzz() got = %v, wantResult %v", result, tt.wantResult)
				return
			}
		})
	}
}

func Test_fnbService_GetFnbRequestsHistory(t *testing.T) {

	now := time.Now()
	type params struct {
		limit *uint64

		dalInteractionsAmount int
		dalErrorOnInteraction error
	}
	tests := []struct {
		name              string
		config            *Config
		params            params
		mockedDALResponse interface{}
		wantResult        []*models.FnbRequest
		wantErr           error
	}{
		{
			name:   "Scenario 1 - OK - 1 element in history",
			config: &Config{Zero: 1},
			params: params{
				limit:                 nil,
				dalInteractionsAmount: 1,
				dalErrorOnInteraction: nil,
			},
			mockedDALResponse: []*models.FnbRequest{
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
			name:   "Scenario 1 - OK - 2 elements in history",
			config: &Config{Zero: 1},
			params: params{
				limit:                 nil,
				dalInteractionsAmount: 1,
				dalErrorOnInteraction: nil,
			},
			mockedDALResponse: []*models.FnbRequest{
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
				{
					ID:          9,
					RequestDate: &now,
					N1:          10,
					S1:          "ten",
					N2:          11,
					S2:          "eleven",
					Limit:       12,
				},
			},
			wantResult: []*models.FnbRequest{
				{
					ID:          5,
					RequestDate: &now,
					N1:          6,
					S1:          "six",
					N2:          7,
					S2:          "seven",
					Limit:       8,
				},
				{
					ID:          9,
					RequestDate: &now,
					N1:          10,
					S1:          "ten",
					N2:          11,
					S2:          "eleven",
					Limit:       12,
				},
			},
			wantErr: nil,
		},
		{
			name:   "Scenario 3 - KO - database returns an error",
			config: &Config{Zero: 1},
			params: params{
				limit:                 nil,
				dalInteractionsAmount: 1,
				dalErrorOnInteraction: errors.New("any error"),
			},
			mockedDALResponse: []*models.FnbRequest{},
			wantResult:        nil,
			wantErr:           fmt.Errorf("couldn't retrieve fizz and buzz requests history: %w", errors.New("any error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			dalMock := mocks.NewMockIDAL(ctrl)

			if tt.params.dalErrorOnInteraction == nil {
				dalMock.EXPECT().GetFnbRequestsHistory(tt.params.limit).Times(tt.params.dalInteractionsAmount).Return(tt.mockedDALResponse, nil)
			} else {
				dalMock.EXPECT().GetFnbRequestsHistory(tt.params.limit).Times(tt.params.dalInteractionsAmount).Return(nil, tt.params.dalErrorOnInteraction)
			}

			service := &fnbService{dal: dalMock, zero: tt.config.Zero}

			got, err := service.GetFnbRequestsHistory(tt.params.limit)
			if !reflect.DeepEqual(err, tt.wantErr) {
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
