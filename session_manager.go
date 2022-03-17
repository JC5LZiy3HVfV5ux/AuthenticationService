package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Реализует интерфес SessionManager из storage.go
type Session struct {
	collection *mongo.Collection
}

func NewSessionManager(storage *MongoStorage) SessionManager {
	session := &Session{}
	session.collection = storage.GetDataBase().Collection("sessions")
	return session
}

func (s *Session) Insert(value SessionDetails) error {
	if _, err := s.collection.InsertOne(context.Background(), value); err != nil {
		return err
	}

	return nil
}

func (s *Session) Get(guid string) (SessionDetails, error) {
	var result SessionDetails

	if err := s.collection.FindOne(context.TODO(), bson.D{{"guid", guid}}).Decode(&result); err != nil {
		return SessionDetails{}, err
	}

	return result, nil
}

func (s *Session) Delete(guid string) error {
	if _, err := s.collection.DeleteOne(context.Background(), bson.D{{"guid", guid}}); err != nil {
		return err
	}

	return nil
}
