package dal

import "time"

type fnbRequest struct {
	tableName   struct{}   `pg:"fnb_requests"`
	ID          uint64     `json:"id" pg:"id"`
	RequestDate *time.Time `json:"-" pg:"request_date"`
	N1          uint32     `json:"n1" pg:"n1"`
	S1          string     `json:"s1" pg:"s1"`
	N2          uint32     `json:"n2" pg:"n2"`
	S2          string     `json:"s2" pg:"s2"`
	Limit       uint64     `json:"limit" pg:"rlimit"` // limit can't be taken, so request limit is abbreviated as rlimit
}

type fnbRequestInputStats struct {
	N1    uint32 `json:"n1" pg:"n1"`
	S1    string `json:"s1" pg:"s1"`
	N2    uint32 `json:"n2" pg:"n2"`
	S2    string `json:"s2" pg:"s2"`
	Limit uint64 `json:"limit" pg:"rlimit"` // limit can't be taken, so request limit is abbreviated as rlimit
	Count uint64 `json:"count" pg:"fnb_count"`
}
