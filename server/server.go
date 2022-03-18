package server

import (
	"log"
	"net/http"

	"authentication/config"
	"authentication/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	AuthService AuthenticationService
	Storage     storage.Storage
	Router      *mux.Router
}

var Srv *Server

func NewServer() {
	Srv = &Server{}

	Srv.AuthService = NewAuthService()
	Srv.Storage = storage.NewMongoStorage()
	Srv.Router = mux.NewRouter()
}

func StartServer() {
	log.Printf("Сервер запущен на %s ... ", config.Conf.Server.ListenAddress)
	log.Printf("Тестовый режим: %t ", config.Conf.TestMode)

	var router = Srv.Router

	if err := http.ListenAndServe(config.Conf.Server.ListenAddress, router); err != nil {
		StopServer()
		log.Fatal("Не удалось запустить сервер: " + err.Error())
	}
}

func StopServer() {
	Srv.Storage.Close()
	log.Print("Сервер остановлен ... ")
}
