package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		StopServer()
		os.Exit(1)
	}()

	NewServer()
	InitApiV1()
	StartServer()
}
