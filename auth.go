package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAuth(r *mux.Router) {
	auth := r.PathPrefix("/auth").Subrouter()

	auth.Use(authMiddleware)
	auth.HandleFunc("/token", token).Methods("GET")
	auth.HandleFunc("/refresh", refresh).Methods("POST")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func token(w http.ResponseWriter, r *http.Request) {
	auth := Srv.AuthService.AuthByToken()
	session := Srv.Storage.Session()

	guid := r.URL.Query().Get("guid")

	if ok := auth.isValidGuid(guid); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := session.Get(guid); err == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := auth.NewAuth(guid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refresh_token_hash, err := auth.HashToken(response.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionDetails, err := NewSessionDetails(guid, refresh_token_hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = session.Insert(sessionDetails); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	output, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func refresh(w http.ResponseWriter, r *http.Request) {
	auth := Srv.AuthService.AuthByToken()
	session := Srv.Storage.Session()

	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok := auth.CompareRefreshAndAccessToken(request.RefreshToken, request.AccessToken); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	access_token, err := auth.ParseAccessToken(request.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok := auth.IsExpiredSAccessToken(access_token.Claims.ExpiresAt); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionDetails, err := session.Get(access_token.Claims.Guid)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ok := auth.CompareHashAndToken(request.RefreshToken, sessionDetails.RefreshToken); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ok := sessionDetails.IsExpiredSessionDetails(); ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := session.Delete(access_token.Claims.Guid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := auth.NewAuth(access_token.Claims.Guid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refresh_token_hash, err := auth.HashToken(response.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionDetails, err = NewSessionDetails(access_token.Claims.Guid, refresh_token_hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = session.Insert(sessionDetails); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	output, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
