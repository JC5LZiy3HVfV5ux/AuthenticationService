package service

type Service interface {
	CreateAuthSession(guid string) (*AuthSession, int, error)
	DeleteAuthSession(accessToken, refreshToken string) (string, int, error)
}
