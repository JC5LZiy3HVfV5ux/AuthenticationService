package utils

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func NewRefreshToken(accessToken string) (string, error) {

	seed := uuid.New().String()
	b64Token := base64.RawURLEncoding.EncodeToString([]byte(seed))

	refreshToken := b64Token + accessToken[len(accessToken)-6:]

	return refreshToken, nil
}
