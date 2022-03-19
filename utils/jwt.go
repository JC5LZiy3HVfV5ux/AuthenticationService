package utils

import (
	"errors"
	"time"

	"authentication/config"

	"github.com/dgrijalva/jwt-go"
)

type AccessToken struct {
	Token  string
	Claims jwt.StandardClaims
}

func NewAccessToken(guid string) (*AccessToken, error) {
	timeDelta, err := time.ParseDuration(config.Conf.AccessTokenTimeDelta)
	if err != nil {
		return &AccessToken{}, err
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(timeDelta).UTC().Unix(),
		Id:        guid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(config.Conf.SecretKey))
	if err != nil {
		return &AccessToken{}, err
	}

	return &AccessToken{
		Token:  accessToken,
		Claims: claims,
	}, nil
}

func ParseAccessToken(tokenString string) (*AccessToken, error) {
	type claims struct {
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.SecretKey), nil
		},
	)

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors != jwt.ValidationErrorExpired {
			return &AccessToken{}, err
		}
	}

	tokenData, ok := token.Claims.(*claims)
	if !ok {
		return &AccessToken{}, errors.New("не удалось извлечь данные из access_token")
	}

	return &AccessToken{
		Token:  tokenString,
		Claims: tokenData.StandardClaims,
	}, nil
}
