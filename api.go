package main

import "github.com/gorilla/mux"

func InitApiV1(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	InitAuth(apiRouter)
}
