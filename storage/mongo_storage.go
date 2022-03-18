package storage

import (
	"context"
	"log"

	"authentication/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Реализует интерфес Storage из storage.go 
type MongoStorage struct {
	client  *mongo.Client
	session SessionManager
}

func NewMongoStorage() Storage {
	storage := &MongoStorage{}

	storage.client = setupConnection()
	storage.session = NewSessionManager(storage)
	return storage
}

func setupConnection() *mongo.Client {
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Conf.DatabaseURI))
	if err != nil {
		log.Fatal("Ошибка при попытке коннекта с MongoDB: " + err.Error())
	}

	if err = db.Ping(context.Background(), nil); err != nil {
		log.Fatal("Ошибка при попытке пинга к MongoDB: " + err.Error())
	}

	return db
}

func (m *MongoStorage) GetDataBase() *mongo.Database {
	return m.client.Database(config.Conf.DatabaseName)
}

func (m *MongoStorage) Session() SessionManager {
	return m.session
}

func (m *MongoStorage) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		log.Fatal("Ошибка при закрытии соединения с MongoDB: " + err.Error())
	}
}
