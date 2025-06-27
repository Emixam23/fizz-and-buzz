package retry

import (
	"math/rand"
	"time"
)

// StopRetry is a wrapped error for stopping the retry loop before all the attempts are made
type StopRetry struct {
	Err error
}

// Error returns the internal error of the stop retry object
func (s *StopRetry) Error() string {
	return s.Err.Error()
}

// Retry is a simple logical function which will retry {attempts} time a function until it doesn't return an error
// If the attempts of retries get reached, then a stop retry error gets returned with the end error internally
func Retry(attempts int, sleep time.Duration, f func(attempts int, sleep time.Duration) error) error {

	if err := f(attempts, sleep); err != nil {
		if s, ok := err.(*StopRetry); ok {
			// Return the original error for later checking
			return s
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}

		return err
	}

	return nil
}
