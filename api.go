package main

import (
	"github.com/gorilla/mux"
)

// API implements an API interface
type API interface {
	Prefix() string
	Route(*mux.Router)
}

func routeAPI(a API, r *mux.Router) {
	sr := r.PathPrefix(a.Prefix()).Subrouter()
	a.Route(sr)
}
