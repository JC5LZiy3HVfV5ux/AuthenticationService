package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Session struct {
	Guid         string `bson:"guid"`
	RefreshToken string `bson:"refresh_token"`
	ExpiredAt    int64  `bson:"expires_at"`
}

type SessionManager struct {
	collection *mongo.Collection
}

func NewSessionManager(db *mongo.Database, collection string) *SessionManager {
	return &SessionManager{
		collection: db.Collection(collection),
	}
}

func (s *SessionManager) Insert(value *Session) error {
	if _, err := s.collection.InsertOne(context.Background(), value); err != nil {
		return err
	}

	return nil
}

func (s *SessionManager) Get(guid string) (*Session, error) {
	var result Session

	if err := s.collection.FindOne(context.TODO(), bson.D{{"guid", guid}}).Decode(&result); err != nil {
		return &Session{}, err
	}

	return &result, nil
}

func (s *SessionManager) Delete(guid string) error {
	if _, err := s.collection.DeleteOne(context.Background(), bson.D{{"guid", guid}}); err != nil {
		return err
	}

	return nil
}
