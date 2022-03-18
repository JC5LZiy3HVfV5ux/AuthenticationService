package storage

import (
	"time"

	"authentication/config"
)

type SessionDetails struct {
	Guid         string `bson:"guid"`
	RefreshToken string `bson:"refresh_token"`
	CreatedAt    int64  `bson:"created_at"`
	ExpiresAt    int64  `bson:"expires_at"`
}

func NewSessionDetails(guid, refreshToken string) (SessionDetails, error) {
	timeDelta, err := time.ParseDuration(config.Conf.RefreshTokenTimeDelta)
	if err != nil {
		return SessionDetails{}, err
	}

	return SessionDetails{
		Guid:         guid,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().UTC().Unix(),
		ExpiresAt:    time.Now().Add(timeDelta).UTC().Unix(),
	}, nil
}

func (s *SessionDetails) IsExpiredSessionDetails() bool {
	now := time.Now().UTC()
	exp := time.Unix(s.ExpiresAt, 0)
	return now.After(exp)
}
