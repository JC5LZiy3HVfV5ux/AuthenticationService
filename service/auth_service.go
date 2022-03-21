package service

import (
	"authentication/config"
	"authentication/storage"
	"authentication/utils"
	"errors"
	"time"
)

type AuthSession struct {
	AccessToken  string
	ExpiresAt    int64
	RefreshToken string
}

type AuthService struct {
	authStorage storage.Manager
}

func NewAuthService(authStorage storage.Manager) *AuthService {
	return &AuthService{
		authStorage: authStorage,
	}
}

func (a *AuthService) CreateAuthSession(guid string) (*AuthSession, int, error) {
	if ok := utils.IsValidGuid(guid); !ok {
		return &AuthSession{}, 400, errors.New("bad request")
	}

	if _, err := a.authStorage.Get(guid); err == nil {
		return &AuthSession{}, 401, errors.New("unauthorized")
	}

	accessToken, err := utils.NewAccessToken(guid)
	if err != nil {
		return &AuthSession{}, 500, errors.New("internal server error")
	}

	refreshToken, err := utils.NewRefreshToken(accessToken.Token)
	if err != nil {
		return &AuthSession{}, 500, errors.New("internal server error")
	}

	refreshTokenHash, err := utils.HashToken(refreshToken)
	if err != nil {
		return &AuthSession{}, 500, errors.New("internal server error")
	}

	timeDelta, err := time.ParseDuration(config.Conf.RefreshTokenTimeDelta)
	if err != nil {
		return &AuthSession{}, 500, errors.New("internal server error")
	}

	if err := a.authStorage.Insert(&storage.Session{
		Guid:         guid,
		RefreshToken: refreshTokenHash,
		ExpiredAt:    time.Now().Add(timeDelta).UTC().Unix(),
	}); err != nil {
		return &AuthSession{}, 500, errors.New("internal server error")
	}

	return &AuthSession{
		AccessToken:  accessToken.Token,
		ExpiresAt:    accessToken.Claims.ExpiresAt,
		RefreshToken: refreshToken,
	}, 201, nil
}

func (a *AuthService) DeleteAuthSession(accessToken, refreshToken string) (string, int, error) {
	if ok := utils.CompareRefreshAndAccessToken(accessToken, refreshToken); !ok {
		return "", 400, errors.New("bad request")
	}

	accessTokenData, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return "", 400, errors.New("bad request")
	}

	if ok := utils.IsExpired(accessTokenData.Claims.ExpiresAt); !ok {
		return "", 401, errors.New("unauthorized")
	}

	session, err := a.authStorage.Get(accessTokenData.Claims.Id)

	if err != nil {
		return "", 401, errors.New("unauthorized")
	}

	if ok := utils.CompareHashAndToken(refreshToken, session.RefreshToken); !ok {
		return "", 401, errors.New("unauthorized")
	}

	if ok := utils.IsExpired(session.ExpiredAt); ok {
		return "", 401, errors.New("unauthorized")
	}

	if err := a.authStorage.Delete(accessTokenData.Claims.Id); err != nil {
		return "", 500, errors.New("internal server error")
	}

	return accessTokenData.Claims.Id, 200, nil
}
