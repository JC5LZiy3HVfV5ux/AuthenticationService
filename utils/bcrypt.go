package utils

import "golang.org/x/crypto/bcrypt"

func HashToken(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}