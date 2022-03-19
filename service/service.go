package service

type Service interface {
	CreateAuthSession(guid string) (*AuthSession, int, error)
	RefreshAuthSession(accessToken, refreshToken string) (string, int, error)
}
