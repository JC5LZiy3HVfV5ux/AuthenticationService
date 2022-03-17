package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	AuthService AuthenticationService
	Storage     Storage
	Router      *mux.Router
}

var Srv *Server

func NewServer() {
	Srv = &Server{}

	Srv.AuthService = NewAuthService()
	Srv.Storage = NewMongoStorage()
	Srv.Router = mux.NewRouter()
}

func StartServer() {
	log.Printf("Сервер запущен на %s ... ", Conf.Server.ListenAddress)
	log.Printf("Тестовый режим: %t ", Conf.TestMode)

	var router = Srv.Router

	if err := http.ListenAndServe(Conf.Server.ListenAddress, router); err != nil {
		StopServer()
		log.Fatal("Не удалось запустить сервер: " + err.Error())
	}
}

func StopServer() {
	Srv.Storage.Close()
	log.Print("Сервер остановлен ... ")
}
