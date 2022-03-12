package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

var Srv *Server

func NewServer() {
	Srv = &Server{}
	Srv.Router = mux.NewRouter()
}

func StartServer() {
	log.Printf("Сервер запущен на %s ... ", Conf.Server.ListenAddress)
	log.Printf("Тестовый режим: %t ", Conf.TestMode)

	var router = Srv.Router

	if err := http.ListenAndServe(Conf.Server.ListenAddress, router); err != nil {
		log.Printf("Не удалось запустить сервер: %s", err)
		StopServer()
		os.Exit(1)
	}
}

func StopServer() {
	log.Print("Сервер остановлен ... ")
}
