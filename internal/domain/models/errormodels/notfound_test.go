package errormodels

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewNotFoundError(t *testing.T) {
	type params struct {
		err error
	}
	tests := []struct {
		name   string
		params params
		want   *NotFoundError
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				err: errors.New("any error"),
			},
			want: &NotFoundError{errors.New("any error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.params.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFoundError_MarshalJSON(t *testing.T) {
	type fields struct {
		error error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Scenario 1 - OK",
			fields: fields{
				error: errors.New("any error"),
			},
			want:    []byte("{\"err\":\"any error\"}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NotFoundError{
				error: tt.fields.error,
			}
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
