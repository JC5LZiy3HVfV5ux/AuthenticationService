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
	log.Print("Сервер запущен на localhost:8080 ... ")

	var router = Srv.Router

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		log.Printf("Не удалось запустить сервер: %s", err)
		StopServer()
		os.Exit(1)
	}
}

func StopServer() {
	log.Print("Сервер остановлен ... ")
}
