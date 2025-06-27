package ui

import "net/http"

// IRestAPI defines a REST API that will expose the IApp
type IRestAPI interface {
	ListenAndServe() error
	Shutdown() error

	// ServeHTTP is for integration testing requirement
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// Config contains the configuration for the IRestAPI
type Config struct {
	Host string
	Port uint32
	Mode string
}
