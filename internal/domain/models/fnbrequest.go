package models

import "time"

// FnbRequest contains the arguments of a Fizz And Buzz command request
type FnbRequest struct {
	ID          uint64
	RequestDate *time.Time
	N1          uint32
	S1          string
	N2          uint32
	S2          string
	Limit       uint64
}
