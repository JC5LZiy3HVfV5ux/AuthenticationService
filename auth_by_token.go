package main

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Реализует интерфес AuthenticationService из auth_service.go
type AuthService struct {
	authService AuthToken
}

// Реализует интерфес AuthToken из auth_service.go
type AuthTokenStruct struct{}

type AccessToken struct {
	Token  string
	Claims Claims
}

type Claims struct {
	Guid string
	jwt.StandardClaims
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"created_at"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"expires_at"`
}

type AuthRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAccessToken(guid string) (*AccessToken, error) {
	timeDelta, err := time.ParseDuration(Conf.AccessTokenTimeDelta)
	if err != nil {
		return &AccessToken{}, err
	}

	claims := Claims{
		Guid: guid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(timeDelta).UTC().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(Conf.SecretKey))
	if err != nil {
		return &AccessToken{}, err
	}

	return &AccessToken{
		Token:  accessToken,
		Claims: claims,
	}, nil
}

func NewRefreshToken(accessToken string) (string, error) {

	seed := uuid.New().String()
	b64Token := base64.RawURLEncoding.EncodeToString([]byte(seed))

	refreshToken := b64Token + accessToken[len(accessToken)-6:]

	return refreshToken, nil
}

func NewAuthService() AuthenticationService {
	auth := &AuthService{}

	auth.authService = NewAuthToken()
	return auth
}

func (a *AuthService) AuthByToken() AuthToken {
	return a.authService
}

func NewAuthToken() AuthToken {
	return &AuthTokenStruct{}
}

func (at *AuthTokenStruct) isValidGuid(guid string) bool {
	if _, err := uuid.Parse(guid); err != nil {
		return false
	}

	return true
}

func (at *AuthTokenStruct) NewAuth(guid string) (*AuthResponse, error) {
	access_token, err := NewAccessToken(guid)
	if err != nil {
		return &AuthResponse{}, err
	}

	refresh_token, err := NewRefreshToken(access_token.Token)
	if err != nil {
		return &AuthResponse{}, err
	}

	return &AuthResponse{
		AccessToken:  access_token.Token,
		TokenType:    "Bearer",
		ExpiresAt:    access_token.Claims.ExpiresAt,
		RefreshToken: refresh_token,
		CreatedAt:    access_token.Claims.IssuedAt}, nil
}

func (at *AuthTokenStruct) CompareRefreshAndAccessToken(refreshToken, accessToken string) bool {
	if len(refreshToken) < 7 || len(accessToken) < 7 {
		return false
	}

	rightAt := accessToken[len(accessToken)-6:]
	rightRt := refreshToken[len(refreshToken)-6:]

	return rightAt == rightRt
}

func (at *AuthTokenStruct) HashToken(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (at *AuthTokenStruct) CompareHashAndToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

func (at *AuthTokenStruct) ParseAccessToken(tokenString string) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Conf.SecretKey), nil
		},
	)

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors != jwt.ValidationErrorExpired {
			return &AccessToken{}, err
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return &AccessToken{}, errors.New("не удалось извлечь данные из access_token")
	}

	return &AccessToken{
		Token:  tokenString,
		Claims: *claims,
	}, nil
}

func (at *AuthTokenStruct) IsExpiredSAccessToken(expiresAt int64) bool {
	now := time.Now().UTC()
	exp := time.Unix(expiresAt, 0)
	return now.After(exp)
}
