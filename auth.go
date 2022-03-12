package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAuth(r *mux.Router) {
	auth := r.PathPrefix("/auth").Subrouter()

	auth.Use(authMiddleware)
	auth.HandleFunc("/token", token)
	auth.HandleFunc("/refresh", refresh)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func token(w http.ResponseWriter, r *http.Request) {

}

func refresh(w http.ResponseWriter, r *http.Request) {

}
