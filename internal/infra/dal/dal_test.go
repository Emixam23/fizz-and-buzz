package dal

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/emixam23/fizz-and-buzz/internal/domain/infra"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	client, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer client.Close()

	type params struct {
		client *sql.DB
	}
	tests := []struct {
		name       string
		params     params
		wantResult infra.IDAL
		wantErr    error
	}{
		{
			name: "Scenario 1 - OK",
			params: params{
				client: client,
			},
			wantResult: &dal{client: client},
			wantErr:    nil,
		},
		{
			name: "Scenario X - KO - provided postgres client is nil",
			params: params{
				client: nil,
			},
			wantResult: nil,
			wantErr:    errors.New("provided postgres client is nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := New(tt.params.client)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("New() got = %v, want %v", got, tt.wantResult)
				return
			}
		})
	}
}
