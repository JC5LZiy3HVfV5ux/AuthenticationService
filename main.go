package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	InitApiV1(router)

	log.Print("Сервер запущен на localhost:8080 ... ")
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Printf("Не удалось запустить сервер: %s", err)
		os.Exit(1)
	}
}
