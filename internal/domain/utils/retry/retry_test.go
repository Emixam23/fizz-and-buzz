package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const sleep = 1 * time.Microsecond

func TestRetry_Should_retry_until_no_error_is_returned_or_all_attempts_have_been_made(t *testing.T) {
	cases := []struct {
		name                  string
		errors                []error
		expectedAttempts      int
		expectedErrorReturned error
	}{
		{
			name: "Scenario 1 - OK - Should not retry if there is no error",
			errors: []error{
				nil,
			},
			expectedAttempts: 1,
		},
		{
			name: "Scenario 2 - OK - Should stop retrying when no error is returned and return no error",
			errors: []error{
				errors.New("first-error"),
				nil,
				errors.New("last-error"),
			},
			expectedAttempts: 2,
		},
		{
			name: "Scenario 3 - OK - Should retry all attempts if there is always an error and return the last error",
			errors: []error{
				errors.New("first-error"),
				errors.New("second-error"),
				errors.New("last-error"),
			},
			expectedAttempts:      3,
			expectedErrorReturned: errors.New("last-error"),
		},
		{
			name: "Scenario 4 - OK - Should retry all attempts until a stop retry is returned",
			errors: []error{
				errors.New("first-error"),
				&StopRetry{Err: errors.New("second-error-as-stop-retry-error")},
			},
			expectedAttempts:      2,
			expectedErrorReturned: &StopRetry{Err: errors.New("second-error-as-stop-retry-error")},
		},
	}

	for _, tc := range cases {
		attempt := 0
		errorReturned := Retry(len(tc.errors), sleep, func(remainingAttempts int, sleep time.Duration) error {
			attempt++
			return tc.errors[attempt-1]
		})

		assert.Equal(t, tc.expectedAttempts, attempt)
		assert.Equal(t, tc.expectedErrorReturned, errorReturned)
	}
}

func TestRetry_Should_wait_between_each_attempt(t *testing.T) {
	cases := []struct {
		name            string
		attempts        int
		sleep           time.Duration
		minExpectedTime time.Duration
	}{
		{
			name:     "Scenario 1 - OK - Should not wait if there is only one attempt",
			attempts: 1,
			sleep:    1 * time.Second,
		},
		{
			name:            "Scenario 1 - OK - Should wait more and more between each attempt",
			attempts:        5,
			sleep:           1 * time.Millisecond,
			minExpectedTime: 5 * time.Millisecond,
		},
	}

	for _, tc := range cases {
		start := time.Now()
		Retry(tc.attempts, tc.sleep, func(remainingAttempts int, sleep time.Duration) error {
			return errors.New("an-error")
		})
		elapsed := time.Since(start)

		assert.LessOrEqual(t, tc.minExpectedTime.Milliseconds(), elapsed.Milliseconds())
	}
}

func TestStopRetry_Error(t *testing.T) {
	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Scenario 1 - OK",
			fields: fields{
				Err: errors.New("any error"),
			},
			want: "any error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StopRetry{
				Err: tt.fields.Err,
			}
			if got := s.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
