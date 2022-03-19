package utils

import (
	"time"

	"github.com/google/uuid"
)

func IsValidGuid(guid string) bool {
	if _, err := uuid.Parse(guid); err != nil {
		return false
	}

	return true
}

func CompareRefreshAndAccessToken(accessToken, refreshToken string) bool {
	if len(accessToken) < 7 || len(refreshToken) < 7 {
		return false
	}

	rightAt := accessToken[len(accessToken)-6:]
	rightRt := refreshToken[len(refreshToken)-6:]

	return rightAt == rightRt
}

func IsExpired(expiresAt int64) bool {
	now := time.Now().UTC()
	exp := time.Unix(expiresAt, 0)
	return now.After(exp)
}