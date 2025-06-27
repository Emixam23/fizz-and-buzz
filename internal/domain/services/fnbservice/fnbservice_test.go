package fnbservice

import (
	"errors"
	"github.com/golang/mock/gomock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"gitlab.com/emixam23/fizz-and-buzz/tests/mocks"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	ctrl := gomock.NewController(t)
	dalMock := mocks.NewMockIDAL(ctrl)

	type params struct {
		config *Config
		dal    infra.IDAL
	}
	tests := []struct {
		name       string
		params     params
		wantResult IFnbService
		wantErr    error
	}{
		{
			name:       "Scenario 1 - OK",
			params:     params{config: &Config{Zero: 5}, dal: dalMock},
			wantResult: &fnbService{zero: 5, dal: dalMock},
			wantErr:    nil,
		},
		{
			name:       "Scenario 1 - KO - nil config provided",
			params:     params{config: nil, dal: dalMock},
			wantResult: nil,
			wantErr:    errors.New("no config provided"),
		},
		{
			name:       "Scenario 1 - KO - nil dal provided",
			params:     params{config: &Config{Zero: 5}, dal: nil},
			wantResult: nil,
			wantErr:    errors.New("provided database access layer does not seem initialized (nil)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.params.config, tt.params.dal)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("New() got = %v, wantResult %v", got, tt.wantResult)
				return
			}
		})
	}
}
