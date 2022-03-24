package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"authentication/config"
	"authentication/service"
	"authentication/storage"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	httpServer  *http.Server
	authService service.Service
}

func NewServer() *Server {
	db := setupConnection()

	sessionManager := storage.NewSessionManager(db, "sessions")
	return &Server{
		authService: service.NewAuthService(sessionManager),
	}
}

func (s *Server) StartServer(listenAddress string) error {

	router := mux.NewRouter()

	log.Printf("Сервер запущен на %s ... ", listenAddress)

	RegisterHandlers(router, s.authService)

	s.httpServer = &http.Server{
		Addr:    listenAddress,
		Handler: router,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal("Не удалось запустить сервер: " + err.Error())
		}
	}()

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}

func setupConnection() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Conf.DatabaseURI))
	if err != nil {
		log.Fatalf("Ошибка при попытке создать клиент mongoDB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal("Ошибка при попытке коннекта с MongoDB: " + err.Error())
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Ошибка при попытке пинга к MongoDB: " + err.Error())
	}

	return client.Database(config.Conf.DatabaseName)
}
