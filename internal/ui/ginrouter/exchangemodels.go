package ginrouter

import "time"

/* HEALTH */

type healthStatus struct {
	Status string `json:"status"`
}

/* ERRORS */

type errorStatus struct {
	Reason string `json:"reason"`
}

func newErrorStatus(err error) *errorStatus {
	return &errorStatus{err.Error()}
}

/* EXCHANGE (params) */

type fizzAndBuzzRequestParams struct {
	N1    uint32 `form:"n1" json:"n1" binding:"required,gte=1,nefield=N2"`
	S1    string `form:"s1" json:"s1" binding:"required"`
	N2    uint32 `form:"n2" json:"n2" binding:"required,gte=1,nefield=N1"`
	S2    string `form:"s2" json:"s2" binding:"required"`
	Limit uint64 `form:"limit" json:"limit" binding:"required,gte=1"`
}

type historyRequestParams struct {
	Limit *uint64 `form:"limit"`
}

type statsRequestParams struct {
	Sorted bool `form:"sorted"`
}

/* EXCHANGE (response models) */

type fnbRequest struct {
	ID          uint64     `json:"id"`
	RequestDate *time.Time `json:"request_date"`
	N1          uint32     `json:"n1"`
	S1          string     `json:"s1"`
	N2          uint32     `json:"n2"`
	S2          string     `json:"s2"`
	Limit       uint64     `json:"limit"`
}

type fnbRequestInputStats struct {
	N1    uint32 `json:"n1"`
	S1    string `json:"s1"`
	N2    uint32 `json:"n2"`
	S2    string `json:"s2"`
	Limit uint64 `json:"limit"`
	Count uint64 `json:"count"`
}

/* EXCHANGE (response containers) */

type fizzAndBuzzResponse struct {
	Request fizzAndBuzzRequestParams `json:"request"`
	Result  []string                 `json:"result"`
}
