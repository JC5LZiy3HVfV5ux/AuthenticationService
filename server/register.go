package server

import (
	"authentication/service"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router, authService service.Service) {
	h := NewHandler(authService)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	auth := apiRouter.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/token", h.Token).Methods("GET")
	auth.HandleFunc("/refresh", h.Refresh).Methods("POST")
}
