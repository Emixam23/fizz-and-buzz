package statsservice

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"reflect"
	"testing"
)

func Test_statsService_GetFnbRequestsInputsStats(t *testing.T) {

	type params struct {
		sorted bool
	}
	tests := []struct {
		name              string
		params            params
		mockedDALResponse []*models.FnbRequestInputStats
		mockedDALError    error
		wantResult        []*models.FnbRequestInputStats
		wantErr           error
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				sorted: false,
			},
			mockedDALResponse: []*models.FnbRequestInputStats{
				{
					N1:    1,
					S1:    "one",
					N2:    2,
					S2:    "two",
					Limit: 3,
					Count: 4,
				},
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 8,
				},
			},
			mockedDALError: nil,
			wantResult: []*models.FnbRequestInputStats{
				{
					N1:    1,
					S1:    "one",
					N2:    2,
					S2:    "two",
					Limit: 3,
					Count: 4,
				},
				{
					N1:    5,
					S1:    "five",
					N2:    6,
					S2:    "six",
					Limit: 7,
					Count: 8,
				},
			},
			wantErr: nil,
		},
		{
			name: "Scenario 2 - KO - database returns an error",
			params: params{
				sorted: false,
			},
			mockedDALResponse: nil,
			mockedDALError:    errors.New("any error"),
			wantResult:        nil,
			wantErr:           fmt.Errorf("couldn't retrieve data to compute for statistics: %w", errors.New("any error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			dalMock := mocks.NewMockIDAL(ctrl)

			dalMock.EXPECT().GetFnbRequestsInputsStats(tt.params.sorted).Times(1).Return(tt.mockedDALResponse, tt.mockedDALError)

			service := &statsService{dal: dalMock}

			got, err := service.GetFnbRequestsInputsStats(tt.params.sorted)
			if !reflect.DeepEqual(err, tt.wantErr) {
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

func Test_statsService_GetFnbRequestsMostUsedCombination(t *testing.T) {

	tests := []struct {
		name              string
		mockedDALResponse *models.FnbRequestInputStats
		mockedDALError    error
		wantResult        *models.FnbRequestInputStats
		wantErr           error
	}{
		{
			name: "Scenario 1 - OK",
			mockedDALResponse: &models.FnbRequestInputStats{
				N1:    1,
				S1:    "one",
				N2:    2,
				S2:    "two",
				Limit: 3,
				Count: 4,
			},
			mockedDALError: nil,
			wantResult: &models.FnbRequestInputStats{
				N1:    1,
				S1:    "one",
				N2:    2,
				S2:    "two",
				Limit: 3,
				Count: 4,
			},
			wantErr: nil,
		},
		{
			name:              "Scenario 3 - KO - database returns an error",
			mockedDALResponse: nil,
			mockedDALError:    errors.New("any error"),
			wantResult:        nil,
			wantErr:           fmt.Errorf("couldn't get most used combination: %w", errors.New("any error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			dalMock := mocks.NewMockIDAL(ctrl)

			dalMock.EXPECT().GetFnbRequestsMostUsedCombination().Times(1).Return(tt.mockedDALResponse, tt.mockedDALError)

			service := &statsService{dal: dalMock}

			got, err := service.GetFnbRequestsMostUsedCombination()
			if !reflect.DeepEqual(err, tt.wantErr) {
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
