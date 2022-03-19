package main

import (
	"authentication/config"
	"authentication/server"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

func main() {
	config.NewConfig()

	app := server.NewServer()

	if err := app.StartServer(config.Conf.Server.ListenAddress); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
