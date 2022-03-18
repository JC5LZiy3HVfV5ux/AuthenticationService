package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"authentication/config"
	"authentication/server"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		server.StopServer()
		os.Exit(1)
	}()

	config.NewConfig()
	server.NewServer()
	server.InitApiV1()
	server.StartServer()
}
