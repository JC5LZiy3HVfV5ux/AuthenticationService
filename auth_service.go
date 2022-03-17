package main

type AuthenticationService interface {
	AuthByToken() AuthToken
}

type AuthToken interface {
	isValidGuid(guid string) bool
	NewAuth(guid string) (*AuthResponse, error)
	CompareRefreshAndAccessToken(refreshToken, accessToken string) bool
	HashToken(token string) (string, error)
	CompareHashAndToken(token, hash string) bool
	ParseAccessToken(token string) (*AccessToken, error)
	IsExpiredSAccessToken(expiresAt int64) bool
}
