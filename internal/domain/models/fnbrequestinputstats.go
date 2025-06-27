package models

// FnbRequestInputStats contains statistics regarding the usage of the API
type FnbRequestInputStats struct {
	N1    uint32
	S1    string
	N2    uint32
	S2    string
	Limit uint64
	Count uint64
}
