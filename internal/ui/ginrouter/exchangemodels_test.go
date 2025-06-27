package ginrouter

import (
	"errors"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/models/errormodels"
	"reflect"
	"testing"
)

func Test_newErrorStatus(t *testing.T) {
	type params struct {
		err error
	}
	tests := []struct {
		name   string
		params params
		want   *errorStatus
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				err: errors.New("any error"),
			},
			want: &errorStatus{"any error"},
		},
		{
			name: "Scenario 2 - OK",
			params: params{
				err: errormodels.NewUnprocessableError(errors.New("any error")),
			},
			want: &errorStatus{"any error"},
		},
		{
			name: "Scenario 3 - OK",
			params: params{
				err: errormodels.NewNotFoundError(errors.New("any error")),
			},
			want: &errorStatus{"any error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newErrorStatus(tt.params.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newErrorStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
